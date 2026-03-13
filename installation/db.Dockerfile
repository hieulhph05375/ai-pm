FROM postgres:16-alpine

# Set default environment variables
ENV POSTGRES_USER=devops
ENV POSTGRES_PASSWORD=rbNyb3cTNLxeM
ENV POSTGRES_DB=pm_db

# Copy initialization scripts
COPY installation/docker_postgres_init.sh /docker-entrypoint-initdb.d/init.sh
COPY installation/seed_production.sql /docker-entrypoint-initdb.d/seed_production.sql
COPY backend/internal/infrastructure/db/migrations /docker-entrypoint-initdb.d/migrations

# Expose standard PostgreSQL port
EXPOSE 5432
