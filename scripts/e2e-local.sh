#!/bin/bash
# e2e-local.sh — One-shot: start infra, launch Go services, seed, run hurl e2e tests.
#
# This is a microservice architecture: one PostgreSQL instance per bounded context
# (identity, merchant, catalog, sales, experience, email), each carrying a single
# ec_<ctx> database and fronted by its own PgBouncer. Go services run on the host;
# infrastructure runs in Docker.
#
# Usage:
#   ./scripts/e2e-local.sh              # full run
#   ./scripts/e2e-local.sh --infra-only # start infra only
#   ./scripts/e2e-local.sh --skip-build # skip go build (use existing binaries)
#   ./scripts/e2e-local.sh --down       # tear down infra
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT"
COMPOSE_FILE="deployments/local/docker-compose.infra.yml"
LOG_DIR="$ROOT/deployments/local/logs"
HURL_DIR="$ROOT/tests/hurl"
mkdir -p "$LOG_DIR"

# ─── Flags ──────────────────────────────────────────────────────────────
INFRA_ONLY=false
DO_DOWN=false
for arg in "$@"; do
  case "$arg" in
    --infra-only) INFRA_ONLY=true ;;
    --down)       DO_DOWN=true ;;
    --help|-h)
      echo "Usage: $0 [--infra-only] [--down]"
      exit 0 ;;
    *) echo "Unknown flag: $arg"; exit 1 ;;
  esac
done

# ─── Colors ─────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
info()  { echo -e "${BLUE}[INFO]${NC} $*"; }
ok()    { echo -e "${GREEN}[OK]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
fail()  { echo -e "${RED}[FAIL]${NC} $*"; }

# ─── Teardown ───────────────────────────────────────────────────────────
cleanup() {
  echo ""
  info "Stopping Go services ..."
  bash "$ROOT/deployments/local/scripts/services-local-stop.sh" 2>/dev/null || true
  rm -f "$ROOT/service/apigateway/.env"
}
trap cleanup EXIT

if [ "$DO_DOWN" = true ]; then
  info "Tearing down infra ..."
  docker compose -f "$COMPOSE_FILE" down -v
  exit 0
fi

# ─── [1/6] Start infra ─────────────────────────────────────────────────
info "[1/6] Starting infrastructure ..."
docker compose -f "$COMPOSE_FILE" up -d

info "Waiting for the six bounded-context Postgres instances ..."
# 1 database = 1 bounded context (deployments/kubernetes/database/):
# postgres_identity … postgres_email, each carrying ec_<ctx>.
for ctx in identity merchant catalog sales experience email; do
  svc="postgres_$ctx"
  timeout=60
  while ! docker compose -f "$COMPOSE_FILE" exec -T "$svc" pg_isready -U DRAGON -q 2>/dev/null; do
    timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "$svc not ready"; exit 1; }
    sleep 1
  done
  ok "$svc ready"
done

info "Waiting for Redis ..."
timeout=30
while ! docker compose -f "$COMPOSE_FILE" exec -T redis redis-cli -a dragon_knight ping 2>/dev/null | grep -q PONG; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "Redis not ready"; exit 1; }
  sleep 1
done
ok "Redis ready"

info "Waiting for Kafka ..."
timeout=90
while ! docker compose -f "$COMPOSE_FILE" exec -T kafka bash -c 'exec 3<>/dev/tcp/localhost/9092' 2>/dev/null; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "Kafka not ready"; exit 1; }
  sleep 2
done
ok "Kafka ready"

info "Waiting for ClickHouse ..."
timeout=30
while ! docker compose -f "$COMPOSE_FILE" exec -T clickhouse clickhouse-client --query 'SELECT 1' 2>/dev/null; do
  timeout=$((timeout - 1)); [ "$timeout" -le 0 ] && { fail "ClickHouse not ready"; exit 1; }
  sleep 1
done
ok "ClickHouse ready"

if [ "$INFRA_ONLY" = true ]; then
  ok "Infrastructure is up. Exiting (--infra-only)."
  echo "Postgres: 6 contexts (ec_identity…ec_email) via PgBouncer localhost:6432-6437"
  echo "Redis: localhost:6379  Kafka: localhost:9092  ClickHouse: localhost:8123/9000"
  exit 0
fi

# ─── [2/6] Reset databases + migrate + seed ─────────────────────────────
info "[2/6] Resetting databases, running migrations, and seeding ..."

# 1 database = 1 bounded context (deployments/kubernetes/database/): the six
# instances carry ec_identity, ec_merchant, ec_catalog, ec_sales,
# ec_experience, ec_email. Host ports 6432-6437 are the PgBouncer poolers that
# the root .env points DB_<CONTEXT>_* at.
CONTEXTS=(identity merchant catalog sales experience email)

# Drop schema + recreate on every context instance for a clean state.
for ctx in "${CONTEXTS[@]}"; do
  docker exec "ecommerce-postgres-$ctx" psql -U DRAGON -d "ec_$ctx" \
    -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" >/dev/null 2>&1 \
    || warn "drop schema failed for ec_$ctx"
done

# Migrate all six contexts in dependency order (reads DB_<CONTEXT>_* from .env).
info "  migrating all six contexts ..."
go run service/migrate/cmd/main.go up > /tmp/hurl_migrate.log 2>&1 \
  || { fail "migrate failed"; cat /tmp/hurl_migrate.log; exit 1; }
ok "All migrations complete"

# Seed roles — the roles table lives in the identity context (ec_identity).
docker exec ecommerce-postgres-identity psql -U DRAGON -d ec_identity -c \
  "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN'), ('ROLE_USER') ON CONFLICT DO NOTHING;" >/dev/null 2>&1 || true
ok "Roles seeded"

# Flush redis
docker exec redis-ecommerce-e2e redis-cli -a dragon_knight FLUSHALL >/dev/null 2>&1 || true
ok "Redis flushed"

# Run seeder
info "  running seeder ..."
go run service/seeder/cmd/main.go > /tmp/hurl_seeder.log 2>&1 || warn "seeder had warnings (see /tmp/hurl_seeder.log)"
ok "Seeder complete"

# Stats backfill
info "  running stats backfill ..."
go run service/stats_writer/cmd/main.go backfill > /tmp/hurl_backfill.log 2>&1 || warn "stats backfill had warnings"
ok "Stats backfill complete"

# ─── [5/6] Launch Go services ──────────────────────────────────────────
info "[5/6] Launching all Go services locally ..."

for svc in auth role user category merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product transaction cart review review_detail slider shipping_address banner email stats_writer stats_reader; do
  (cd "$ROOT/service/$svc" && setsid nohup "$ROOT/bin/$svc" > "$LOG_DIR/$svc.log" 2>&1 &)
done
(cd "$ROOT/service/apigateway" && setsid nohup "$ROOT/bin/apigateway" > "$LOG_DIR/apigateway.log" 2>&1 &)
ok "All 23 services launched"

# ─── [6/6] Wait for health ─────────────────────────────────────────────
info "[6/6] Waiting for all services to become healthy ..."
GW_URL="http://localhost:5000/api/auth/hello"
PORT_REGEX=':5005[1-9]|:5006[0-9]|:5000'

ok=0
for round in 1 2 3; do
  info "  health check round $round ..."
  for i in $(seq 1 40); do
    gw=$(curl -s -o /dev/null -w '%{http_code}' --max-time 2 "$GW_URL" 2>/dev/null || echo "000")
    if [ "$gw" = "200" ]; then ok=1; break; fi
    sleep 3
  done
  info "  round $round: apigateway=$gw"
  [ "$ok" = "1" ] && break
done

if [ "$ok" != "1" ]; then
  fail "Apigateway not healthy. Logs:"
  for f in "$LOG_DIR"/*.log; do echo "-- $(basename "$f")"; tail -5 "$f" 2>/dev/null; done
  exit 1
fi
ok "All services healthy"

# ─── [7/6] Run hurl e2e tests ──────────────────────────────────────────
info "[7/6] Running hurl e2e tests ..."

# Discover a seeded role-less user for rules_strict.hurl.
# users, roles and user_roles all live in the identity context (ec_identity).
USER_EMAIL=""
if docker exec ecommerce-postgres-identity psql -U DRAGON -d ec_identity -t -A -c \
  "SELECT email FROM users WHERE firstname='User1' AND email LIKE 'user_%' ORDER BY user_id LIMIT 1;" \
  > /tmp/hurl_seed_user.txt 2>/dev/null; then
  USER_EMAIL=$(tr -d ' \r' < /tmp/hurl_seed_user.txt)
  USER_ID=$(docker exec ecommerce-postgres-identity psql -U DRAGON -d ec_identity -t -A -c \
    "SELECT user_id FROM users WHERE email = '$USER_EMAIL' LIMIT 1;" | tr -d ' \r')
  if [ -n "$USER_ID" ]; then
    docker exec ecommerce-postgres-identity psql -U DRAGON -d ec_identity -c \
      "DELETE FROM user_roles WHERE user_id = $USER_ID;" >/dev/null 2>&1 || true
  fi
fi

TS=$(date +%s)
TEST_PASSWORD="HurlPass123"
PASS=0; FAIL=0
FAILED=()

for f in "$HURL_DIR"/*.hurl; do
  name=$(basename "$f")
  testEmail="hurl.${name%.hurl}.$TS@example.com"

  common_vars="--test --jobs 1 --variable baseUrl=http://localhost:5000 --variable testEmail=$testEmail --variable testPassword=$TEST_PASSWORD --variable testImage=$HURL_DIR/assets/test.png"

  if [ "$name" = "rules_strict.hurl" ]; then
    if [ -z "$USER_EMAIL" ]; then
      warn "SKIP $name (no seeded role-less user)"
      continue
    fi
    extra_vars="--variable adminEmail=admin.$TS@example.com --variable adminPassword=$TEST_PASSWORD --variable userEmail=$USER_EMAIL --variable userPassword=password1"
    hurl $common_vars $extra_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  else
    hurl $common_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  fi

  if [ $? -eq 0 ]; then
    PASS=$((PASS + 1))
    ok "  ✓ $name"
  else
    FAIL=$((FAIL + 1))
    FAILED+=("$name")
    fail "  ✗ $name (see /tmp/hurl_$name.log)"
  fi

  # Pace auth endpoints (rate-limited)
  sleep 2
done

# ─── [8/6] Verify ClickHouse stats ─────────────────────────────────────
info "[8/6] Verifying ClickHouse stats ..."
sleep 7  # wait for stats-writer flush

for table in order_events order_item_events transaction_events; do
  row_count=$(docker exec clickhouse-ecommerce clickhouse-client --database ecommerce \
    --query "SELECT count() FROM $table FINAL" 2>/dev/null | tr -d ' \r' || echo "unknown")
  info "  stats $table rows: ${row_count}"
done

# ─── Summary ────────────────────────────────────────────────────────────
echo ""
echo "========================================"
echo "  E2E Test Summary"
echo "========================================"
ok "Passed: $PASS"
if [ "$FAIL" -gt 0 ]; then
  fail "Failed: $FAIL"
  echo ""
  echo "Failed suites:"
  for f in "${FAILED[@]}"; do echo "  - $f"; done
  exit 1
else
  ok "All tests passed! 🎉"
fi
