#!/bin/bash
set -e

echo "=============================================="
echo "    AIVOVAN QA Automation Suite Execution     "
echo "=============================================="

# Define directories
ROOT_DIR=$(pwd)
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
RESULTS_DIR="$ROOT_DIR/test-results"

# 1. Clean and Setup Results Directory
echo "[1/4] Preparing test results directory..."
rm -rf "$RESULTS_DIR"
mkdir -p "$RESULTS_DIR"

# Ensure test DB is clean and migrations run
echo "[2/4] Setting up backend test database..."
cd "$BACKEND_DIR"
go test ./tests -run ^TestMain$ > /dev/null 2>&1 || true # Initialize DB and migrations

# 2. Run Backend API Tests
echo "[3/4] Running Backend API Tests (Go)..."
cd "$BACKEND_DIR"
go test -json ./tests/... > "$RESULTS_DIR/backend.json" || true
echo "      Backend tests completed. Output saved."

# 3. Run Frontend E2E Tests
echo "[4/4] Running Frontend E2E Tests (Playwright)..."
cd "$FRONTEND_DIR"
# Force JSON reporter and disable UI output for clean parsing
npx playwright test --reporter=json > "$RESULTS_DIR/frontend.json" || true
echo "      Frontend tests completed. Output saved."

# 4. Generate PDF Report
echo "=============================================="
echo "      Generating Final QA PDF Report          "
echo "=============================================="
cd "$BACKEND_DIR"
go run cmd/qa-report/main.go

echo " "
echo "✅ QA run complete. Report available at: $BACKEND_DIR/qa-report.pdf"
