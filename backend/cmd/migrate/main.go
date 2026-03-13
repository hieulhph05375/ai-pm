package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"project-mgmt/backend/internal/infrastructure/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	migrationsDir := "internal/infrastructure/db/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Error reading migrations directory: %v", err)
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	// In a real system, we'd track applied migrations in a table.
	targetMigrations := []string{
		"000024_create_project_rbac_tables.up.sql",
		"000025_migrate_project_member_roles.up.sql",
		"000026_granular_project_permissions.up.sql",
		"000027_create_snapshots_table.up.sql",
		"000028_add_performance_indexes.up.sql",
		"000029_optimize_wbs_search.up.sql",
		"000030_global_category_refactor.up.sql",
		"000031_add_role_id_to_project_stakeholders.up.sql",
		"000032_wbs_category_refactor.up.sql",
	}

	for _, targetMigration := range targetMigrations {
		found := false
		for _, file := range upFiles {
			if file == targetMigration {
				found = true
				break
			}
		}

		if !found {
			log.Fatalf("Migration file %s not found", targetMigration)
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, targetMigration))
		if err != nil {
			log.Fatalf("Error reading migration file: %v", err)
		}

		fmt.Printf("Applying migration: %s\n", targetMigration)
		_, err = db.ExecContext(context.Background(), string(content))
		if err != nil {
			log.Fatalf("Error executing migration %s: %v", targetMigration, err)
		}
	}

	fmt.Println("✅ Migrations applied successfully!")
}
