#!/usr/bin/env bash
#
# backup.sh — PostgreSQL backup for the local e-commerce stack.
#
# Usage:
#   ./scripts/backup.sh [retention_days]
#
# Produces: deployments/local/backups/ecommerce_<timestamp>/ec_<ctx>.sql.gz
#           — one dump per bounded context (identity, merchant, catalog, sales,
#             experience, email) plus a MANIFEST describing the snapshot.
# Prunes whole snapshots older than RETENTION_DAYS (default 7).
#
# Procedure (Fase 6 checklist 11.3):
#   1. Verify every context Postgres container is healthy.
#   2. pg_dump each ec_<ctx> database with --no-owner --no-privileges so the
#      dump is portable to disposable environments.
#   3. Compress with gzip into the snapshot directory.
#   4. Prune old snapshots so retention is bounded.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
COMPOSE_FILE="${SCRIPT_DIR}/../docker-compose.yml"
BACKUP_DIR="${SCRIPT_DIR}/../backups"
RETENTION_DAYS="${1:-7}"
DB_USER="${POSTGRES_USER:-DRAGON}"

# 1 database = 1 bounded context (pkg/database/names.go + deployments/kubernetes/database/).
# Each instance is fronted by its own PgBouncer, so dumps go straight to the
# Postgres service rather than through the pooler.
CONTEXTS=(identity merchant catalog sales experience email)

mkdir -p "${BACKUP_DIR}"

# `compose ps` exits 0 even when a service is absent/stopped, so check the
# actual container state via inspect instead of relying on its exit code.
for ctx in "${CONTEXTS[@]}"; do
  PG_ID=$(docker compose -f "${COMPOSE_FILE}" ps -q "postgres_$ctx" 2>/dev/null || true)
  PG_STATE=$(docker inspect -f '{{.State.Running}}' "$PG_ID" 2>/dev/null || echo false)
  if [ -z "$PG_ID" ] || [ "$PG_STATE" != "true" ]; then
    echo "ERROR: postgres_$ctx container is not running. Start the stack first (just up)." >&2
    exit 1
  fi
done

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
SNAPSHOT_DIR="${BACKUP_DIR}/ecommerce_${TIMESTAMP}"
mkdir -p "${SNAPSHOT_DIR}"

echo "Backing up ${#CONTEXTS[@]} context databases -> ${SNAPSHOT_DIR}"

for ctx in "${CONTEXTS[@]}"; do
  db="ec_${ctx}"
  out="${SNAPSHOT_DIR}/${db}.sql.gz"

  # pipefail makes a failed pg_dump abort the pipeline; the explicit check also
  # catches a partial snapshot left behind by an earlier failure.
  if ! docker compose -f "${COMPOSE_FILE}" exec -T "postgres_$ctx" \
         pg_dump -U "${DB_USER}" -d "${db}" --no-owner --no-privileges \
         | gzip -9 > "${out}"; then
    echo "ERROR: pg_dump failed for ${db}; removing incomplete snapshot." >&2
    rm -rf "${SNAPSHOT_DIR}"
    exit 1
  fi

  # Validate the archive itself: a truncated gzip would restore nothing. An
  # empty (but valid) dump is accepted — a context whose migrations have not
  # been applied yet is a legitimate thing to snapshot.
  if [ ! -s "${out}" ] || ! gzip -t "${out}" 2>/dev/null; then
    echo "ERROR: ${db} dump is empty or corrupt; removing snapshot." >&2
    rm -rf "${SNAPSHOT_DIR}"
    exit 1
  fi

  echo "  ${db}: $(du -h "${out}" | cut -f1)"
done

# MANIFEST records what the snapshot contains so restore.sh can validate it.
{
  echo "created=${TIMESTAMP}"
  echo "user=${DB_USER}"
  for ctx in "${CONTEXTS[@]}"; do
    echo "context=${ctx} db=ec_${ctx}"
  done
} > "${SNAPSHOT_DIR}/MANIFEST"

echo "Backup complete: ${SNAPSHOT_DIR}"

# Prune old snapshots beyond the retention window.
PRUNE_TS=$(date -d "-${RETENTION_DAYS} days" +%Y%m%d 2>/dev/null || date -v-${RETENTION_DAYS}d +%Y%m%d)
PRUNED=0
for d in "${BACKUP_DIR}"/ecommerce_*/; do
  [ -d "$d" ] || continue
  ddate=$(basename "$d" | sed -E 's/ecommerce_([0-9]{8})_.*/\1/')
  if [ -n "$ddate" ] && [ "$ddate" -lt "$PRUNE_TS" ]; then
    rm -rf "$d"
    PRUNED=$((PRUNED + 1))
  fi
done

echo "Pruned ${PRUNED} snapshot(s) older than ${RETENTION_DAYS} day(s)."
echo "Done."
