#!/bin/bash
# ==============================================================================
# Project Management System - Database Initialization Script
# ==============================================================================
#
# Usage:
#   ./init_database.sh              # Docker mode (default)
#   ./init_database.sh --local      # Local PostgreSQL mode
#
# Requirements:
#   Docker mode: docker must be running and pm_db container must be up
#   Local mode:  psql must be installed and postgres must be running
# ==============================================================================

set -e

# --- Configuration ---
DB_CONTAINER="pm_db_1"
DB_NAME="${DB_NAME:-pm_db}"
DB_USER="${DB_USER:-devops}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_PASSWORD="${DB_PASSWORD:-}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
MIGRATIONS_DIR="$BACKEND_DIR/internal/infrastructure/db/migrations"
SEED_FILE="$SCRIPT_DIR/seed_production.sql"

MODE="docker"
if [[ "$1" == "--local" ]]; then
  MODE="local"
fi

# --- Colors for output ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_step() { echo -e "${BLUE}==>${NC} $1"; }
print_ok()   { echo -e "${GREEN}  ✅ $1${NC}"; }
print_warn() { echo -e "${YELLOW}  ⚠️  $1${NC}"; }
print_err()  { echo -e "${RED}  ❌ $1${NC}"; }

# --- Helper: run psql command ---
run_psql() {
  local sql="$1"
  if [[ "$MODE" == "docker" ]]; then
    docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" -c "$sql"
  else
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "$sql"
  fi
}

run_psql_file() {
  local file="$1"
  if [[ "$MODE" == "docker" ]]; then
    docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" < "$file"
  else
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
  fi
}

# ==============================================================================
# Step 0: Pre-flight checks
# ==============================================================================
print_step "Pre-flight checks..."

if [[ "$MODE" == "docker" ]]; then
  if ! docker info &>/dev/null; then
    print_err "Docker is not running. Please start Docker first."
    exit 1
  fi
  if ! docker ps --format '{{.Names}}' | grep -q "^${DB_CONTAINER}$"; then
    print_err "Container '$DB_CONTAINER' is not running."
    print_warn "Run: docker-compose up -d  (from the installation/ directory)"
    exit 1
  fi
  print_ok "Docker and container '$DB_CONTAINER' are running."
else
  if ! command -v psql &>/dev/null; then
    print_err "psql is not installed. Please install PostgreSQL client tools."
    exit 1
  fi
  print_ok "psql is available."
fi

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
  print_err "Migrations directory not found at: $MIGRATIONS_DIR"
  exit 1
fi

if [[ ! -f "$SEED_FILE" ]]; then
  print_err "Seed file not found at: $SEED_FILE"
  exit 1
fi

# ==============================================================================
# Step 1: Create migrations tracking table
# ==============================================================================
print_step "Ensuring migrations tracking table exists..."
run_psql "CREATE TABLE IF NOT EXISTS schema_migrations (filename VARCHAR(255) PRIMARY KEY, applied_at TIMESTAMPTZ DEFAULT NOW());" > /dev/null
print_ok "Tracking table ready."

# ==============================================================================
# Step 2: Apply migrations
# ==============================================================================
print_step "Applying database migrations..."

MIGRATION_COUNT=0
SKIP_COUNT=0

for migration_file in $(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | sort); do
  filename=$(basename "$migration_file")

  # Check if already applied
  APPLIED=$(run_psql "SELECT COUNT(*) FROM schema_migrations WHERE filename='$filename';" 2>/dev/null | grep -E '^\s+[0-9]+' | tr -d ' ')

  if [[ "$APPLIED" == "1" ]]; then
    SKIP_COUNT=$((SKIP_COUNT + 1))
    continue
  fi

  echo -n "  Applying: $filename ... "
  if run_psql_file "$migration_file" > /dev/null 2>&1; then
    run_psql "INSERT INTO schema_migrations (filename) VALUES ('$filename') ON CONFLICT DO NOTHING;" > /dev/null
    echo -e "${GREEN}done${NC}"
    MIGRATION_COUNT=$((MIGRATION_COUNT + 1))
  else
    echo -e "${RED}FAILED${NC}"
    print_err "Migration '$filename' failed. Aborting."
    print_warn "Fix the migration and re-run this script. Already-applied migrations are safely skipped."
    exit 1
  fi
done

print_ok "Migrations: $MIGRATION_COUNT applied, $SKIP_COUNT already up-to-date."

# ==============================================================================
# Step 3: Apply seed data
# ==============================================================================
print_step "Seeding production data (idempotent)..."
if run_psql_file "$SEED_FILE" > /dev/null; then
  print_ok "Seed data applied successfully."
else
  print_warn "Seed may have partially applied (check for conflicts). This is usually safe."
fi

# ==============================================================================
# Step 4: Seed permissions (roles + permissions matrix)
# ==============================================================================
print_step "Seeding RBAC roles and permissions..."
SEED_PERMISSIONS_BINARY="$BACKEND_DIR/cmd/seed_permissions/main.go"

if [[ -f "$SEED_PERMISSIONS_BINARY" ]]; then
  if command -v go &>/dev/null; then
    if [[ "$MODE" == "docker" ]]; then
      # Inside docker, permissions are seeded by the backend init process
      print_warn "Docker mode: RBAC permissions are seeded automatically on backend startup."
    else
      DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
      DATABASE_URL="$DB_URL" go run "$SEED_PERMISSIONS_BINARY" 2>/dev/null && print_ok "RBAC permissions seeded." || print_warn "Permission seeding skipped (may already be seeded)."
    fi
  else
    print_warn "Go not found. Skipping RBAC permission seeding. Run manually if needed."
  fi
else
  print_warn "Permission seeder not found. Skipping."
fi

# ==============================================================================
# Summary
# ==============================================================================
echo ""
echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}  ✅  Database initialization complete!${NC}"
echo -e "${GREEN}================================================================${NC}"
echo ""
echo "  Mode:        $MODE"
echo "  Database:    $DB_NAME"
echo "  Migrations:  $MIGRATION_COUNT new, $SKIP_COUNT skipped"
echo ""
echo "  Default admin credentials:"
echo "    Email:    admin@admin.com"
echo "    Password: password  ← CHANGE THIS IN PRODUCTION"
echo ""
