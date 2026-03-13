package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Find .env file from the root
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	log.Printf("Loading .env from %s", filepath.Join(dir, ".env"))
	godotenv.Load(filepath.Join(dir, ".env"))

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local dev
		dbURL = "postgres://postgres:Caikeo@1234@localhost:5432/pm_db?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	permissions := []struct {
		Name        string
		Description string
	}{
		{"project:create", "Create new projects"},
		{"project:read", "View project details"},
		{"project:update", "Edit project information"},
		{"project:delete", "Delete projects"},
		{"project:export", "Export project data to Excel/CSV"},
		{"task:create", "Create new tasks"},
		{"task:read", "View tasks"},
		{"task:update", "Update tasks"},
		{"task:delete", "Delete tasks"},
		{"task:export", "Export task data"},
		{"user:create", "Create new users"},
		{"user:read", "View users"},
		{"user:update", "Update users"},
		{"user:delete", "Delete users"},
		{"user:export", "Export users data"},
		{"role:create", "Create new roles"},
		{"role:read", "View roles and permissions"},
		{"role:update", "Update roles and permissions"},
		{"role:delete", "Delete roles"},
		{"timesheet:create", "Create timesheets"},
		{"timesheet:read", "View timesheets"},
		{"timesheet:update", "Update timesheets"},
		{"timesheet:delete", "Delete timesheets"},
		{"timesheet:export", "Export timesheet data"},
		{"category:create", "Create new categories"},
		{"category:read", "View categories"},
		{"category:update", "Update categories"},
		{"category:delete", "Delete categories"},
	}

	for _, p := range permissions {
		var id int
		err := db.QueryRow("INSERT INTO permissions (name, description) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET description = EXCLUDED.description RETURNING id", p.Name, p.Description).Scan(&id)
		if err != nil {
			log.Printf("Error inserting permission %s: %v", p.Name, err)
		} else {
			fmt.Printf("Inserted permission %s with ID %d\n", p.Name, id)
		}
	}

	// Build list of valid names
	var validNames []string
	for _, p := range permissions {
		validNames = append(validNames, "'"+p.Name+"'")
	}
	validNamesQueryStr := strings.Join(validNames, ",")

	// Delete obsolete role_permissions to satisfy FK
	_, err = db.Exec(fmt.Sprintf("DELETE FROM role_permissions WHERE permission_id IN (SELECT id FROM permissions WHERE name NOT IN (%s))", validNamesQueryStr))
	if err != nil {
		log.Printf("Failed to clean obsolete role_permissions: %v", err)
	}

	// Delete obsolete permissions
	res, err := db.Exec(fmt.Sprintf("DELETE FROM permissions WHERE name NOT IN (%s)", validNamesQueryStr))
	if err != nil {
		log.Printf("Failed to clean obsolete permissions: %v", err)
	} else {
		deletedCount, _ := res.RowsAffected()
		if deletedCount > 0 {
			fmt.Printf("Cleaned up %d obsolete permissions from database.\n", deletedCount)
		}
	}

	fmt.Println("Permissions seeded and cleaned successfully!")
}
