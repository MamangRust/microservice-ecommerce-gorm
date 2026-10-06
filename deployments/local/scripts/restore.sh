#!/usr/bin/env bash
#
# restore.sh — PostgreSQL restore for the local e-commerce stack.
#
# Usage:
#   ./scripts/restore.sh [--yes] <snapshot_dir>
#
# The snapshot directory is exactly what backup.sh produces: one ec_<ctx>.sql.gz
# per bounded context. Every context present in the snapshot is restored, so this
# is destructive — a confirmation prompt is shown unless --yes is passed.
#
#   ./scripts/restore.sh --yes backups/ecommerce_20260101_000000
#
# Procedure (Fase 6 checklist 11.3):
#   1. Validate the snapshot directory and each dump.
#   2. Confirm intent (destructive).
#   3. Drop + recreate the public schema in each context database.
#   4. Restore each dump.
#   5. Report success/failure.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${SCRIPT_DIR}/../docker-compose.yml"
DB_USER="${POSTGRES_USER:-DRAGON}"

AUTO_YES=0
SNAPSHOT_DIR=""

for arg in "$@"; do
  case "$arg" in
    --yes) AUTO_YES=1 ;;
    *)     SNAPSHOT_DIR="$arg" ;;
  esac
done

if [ -z "$SNAPSHOT_DIR" ]; then
  echo "Usage: $0 [--yes] <snapshot_dir>" >&2
  echo "  snapshot_dir: a directory produced by backup.sh (contains ec_<ctx>.sql.gz)" >&2
  exit 1
fi

if [ ! -d "$SNAPSHOT_DIR" ]; then
  echo "ERROR: snapshot directory not found: ${SNAPSHOT_DIR}" >&2
  exit 1
fi

# Validate every dump BEFORE the destructive drop: a truncated/empty/corrupt file
# would wipe a context database and restore nothing.
DUMPS=()
for dump in "${SNAPSHOT_DIR}"/ec_*.sql.gz; do
  [ -e "$dump" ] || continue
  if [ ! -s "$dump" ] || ! gzip -t "$dump" 2>/dev/null; then
    echo "ERROR: $(basename "$dump") is empty or not a valid gzip archive; refusing to restore." >&2
    exit 1
  fi
  DUMPS+=("$dump")
done

if [ "${#DUMPS[@]}" -eq 0 ]; then
  echo "ERROR: no ec_<ctx>.sql.gz dumps found in ${SNAPSHOT_DIR}" >&2
  exit 1
fi

# Every dump's context must have a running Postgres instance to restore into.
for dump in "${DUMPS[@]}"; do
  ctx="$(basename "$dump" .sql.gz)"
  ctx="${ctx#ec_}"
  PG_ID=$(docker compose -f "${COMPOSE_FILE}" ps -q "postgres_$ctx" 2>/dev/null || true)
  PG_STATE=$(docker inspect -f '{{.State.Running}}' "$PG_ID" 2>/dev/null || echo false)
  if [ -z "$PG_ID" ] || [ "$PG_STATE" != "true" ]; then
    echo "ERROR: postgres_$ctx container is not running. Start the stack first (just up)." >&2
    exit 1
  fi
done

if [ "$AUTO_YES" -ne 1 ]; then
  echo "Snapshot contains ${#DUMPS[@]} context database(s):"
  for dump in "${DUMPS[@]}"; do echo "  - $(basename "$dump" .sql.gz)"; done
  read -r -p "Restore will DROP and recreate the public schema in each of them. Continue? [y/N] " answer
  case "$answer" in
    y|Y) ;;
    *) echo "Aborted."; exit 1 ;;
  esac
fi

echo "Restoring ${#DUMPS[@]} context database(s) from ${SNAPSHOT_DIR}"

for dump in "${DUMPS[@]}"; do
  db="$(basename "$dump" .sql.gz)"
  ctx="${db#ec_}"
  container="postgres_$ctx"

  echo "  ${db} ..."

  # Terminate other backends first: PgBouncer (and any running Go service) holds
  # sessions that would block DROP SCHEMA. Dropping the schema rather than the
  # database keeps the pooler's target DB intact.
  if ! docker compose -f "${COMPOSE_FILE}" exec -T "$container" \
         psql -U "${DB_USER}" -d "${db}" -v ON_ERROR_STOP=1 \
         -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = current_database() AND pid <> pg_backend_pid();" \
         -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" >/dev/null; then
    echo "ERROR: failed to reset schema for ${db}" >&2
    exit 1
  fi

  if ! gunzip -c "$dump" | docker compose -f "${COMPOSE_FILE}" exec -T "$container" \
         psql -U "${DB_USER}" -d "${db}" -v ON_ERROR_STOP=1 >/dev/null; then
    echo "ERROR: failed to restore ${db} from $(basename "$dump")" >&2
    exit 1
  fi

  echo "    restored"
done

echo "Restore complete: ${#DUMPS[@]} context database(s) back in service."
