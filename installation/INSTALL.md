# Application Installation Guide (DevOps)

This guide provides instructions for deploying the Project Management System in a production environment using Docker.

## Prerequisites

- Docker (v24.0+)
- Docker Compose (v2.20+)
- Linux/Amd64 or Arm64 environment
- Minimum 2GB RAM, 1vCPU

## 1. Environment Configuration

Clone the repository, enter the installation directory, and create your production environment file from the template:

```bash
cd installation
cp .env.example .env
```

Edit `.env` and MUST update the following security-critical variables:

- `DB_PASSWORD`: Set a strong password for the PostgreSQL database.
- `JWT_SECRET`: Generate a random string (minimum 32 characters) used for signing auth tokens.

## 2. Deployment

The system is orchestrated using Docker Compose. Execute high-level build and run:

```bash
# Run from the installation/ directory
docker-compose up -d --build
```

### Services Started

| Service | Image | Internal Port | External Port | Role |
|---------|-------|---------------|---------------|------|
| `pm_db` | postgres:16-alpine | 5432 | 5432 | Primary SQL Storage |
| `pm_backend` | custom (go) | 8081 | 8081 | Core API Engine |
| `pm_frontend` | custom (nginx)| 80 | 80 | UI & Reverse Proxy |

## 3. Database Initialization

Run the init script **once** after the containers are up. It applies all migrations and seeds the production data:

```bash
# Run from the installation/ directory
./init_database.sh
```

This script will:

1. ✅ Apply all database schema migrations (idempotent — safe to re-run)
2. ✅ Seed all system categories (Project Status, WBS Node Types, Holiday Types, etc.)
3. ✅ Create the default admin user

**Local development (without Docker):**

```bash
# Set DB credentials first
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=yourpassword DB_NAME=pm_db

./init_database.sh --local
```

### Default Login

After initialization, log in with:

| Field    | Value              |
|----------|--------------------|
| Email    | `admin@admin.com`  |
| Password | `password`         |

> ⚠️ **Change the default password immediately** in the Users settings page.

## 4. Verification

Verify the system is operational:

1. **Frontend**: Open `http://your-server-ip` (should show Login page).
2. **API Health**: `curl http://localhost:8081/api/v1/health` (should return 200 OK).
3. **Logs**: `docker-compose logs -f backend` to monitor for initialization errors.

## 5. Persistence & Backups

Data is persisted in the following Docker volumes:

- `postgres_data`: DB records and schema.

**Manual Backup:**

```bash
docker exec pm_db pg_dump -U postgres pm_db > backup_$(date +%Y%m%d).sql
```

**Restore:**

```bash
docker exec -i pm_db psql -U postgres pm_db < backup_YYYYMMDD.sql
```

## 6. Troubleshooting

### Port Conflicts

If port 80 or 8081 is already in use by another service on the host, modify the `ports` section in `docker-compose.yml`:

```yaml
# Example: change host port 80 to 8080
frontend:
  ports:
    - "8080:80"
```

### Re-running Initialization

The init script is fully **idempotent** — you can run it multiple times safely. It skips already-applied migrations and uses `ON CONFLICT DO NOTHING` for all seed data.

```bash
./init_database.sh   # run again any time
```

### Connection Issues

Ensure the backend can reach the DB. Inside Docker, the backend refers to the database by the container name `db`.
Verify with: `docker exec -it pm_backend nslookup db`.
