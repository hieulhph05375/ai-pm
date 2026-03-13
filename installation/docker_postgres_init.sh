#!/bin/bash
set -e

echo "Starting automatic Database Initialization for pm_db..."

# Create migrations tracking table
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE IF NOT EXISTS schema_migrations (
        filename VARCHAR(255) PRIMARY KEY, 
        applied_at TIMESTAMPTZ DEFAULT NOW()
    );
EOSQL

echo "Applying migrations..."
MIGRATIONS_DIR="/docker-entrypoint-initdb.d/migrations"

if [ -d "$MIGRATIONS_DIR" ]; then
    for migration_file in $(ls "$MIGRATIONS_DIR"/*.up.sql 2>/dev/null | sort); do
        filename=$(basename "$migration_file")
        
        # Check if already applied
        APPLIED=$(psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -tAc "SELECT COUNT(*) FROM schema_migrations WHERE filename='$filename';")
        
        if [ "$APPLIED" = "1" ]; then
            continue
        fi
        
        echo "  Applying: $filename"
        psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration_file"
        psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "INSERT INTO schema_migrations (filename) VALUES ('$filename') ON CONFLICT DO NOTHING;"
    done
else
    echo "Warning: Migrations directory not found at $MIGRATIONS_DIR"
fi

echo "Applying seed data..."
SEED_FILE="/docker-entrypoint-initdb.d/seed_production.sql"
if [ -f "$SEED_FILE" ]; then
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$SEED_FILE"
    echo "Seed data applied."
else
    echo "Warning: Seed file not found at $SEED_FILE"
fi

echo "Database Initialization complete!"
