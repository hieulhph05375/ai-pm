package tests

import (
	"database/sql"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/infrastructure/config"
	"project-mgmt/backend/internal/infrastructure/db"
	"project-mgmt/backend/internal/presentation"

	_ "github.com/lib/pq"
)

var (
	testDB     *sql.DB
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	// 1. Load test configuration
	// We are running from backend/tests, so we need to go up to root for certs if needed,
	// but let's assume relative paths from root if we set CWD.
	// For TestMain, we'll try to find .env.test
	os.Setenv("DATABASE_URL", "postgres://postgres:Caikeo@1234@localhost:5432/projectmgmt_test?sslmode=disable")
	os.Setenv("JWT_PRIVATE_KEY_PATH", "../certs/private.pem")
	os.Setenv("JWT_PUBLIC_KEY_PATH", "../certs/public.pem")
	os.Setenv("JWT_SECRET", "test-secret-key-2026")

	cfg := config.Load()

	// 2. Setup Test Database
	var err error
	testDB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}
	defer testDB.Close()

	// 3. Clear and Migrate Database
	if err := runMigrations(testDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 4. Seed Fixtures
	if err := SeedAll(testDB); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// 5. Initialize Services and Router
	database, _ := db.NewConnection(cfg) // Wrap the *sql.DB if needed or just use database/sql

	userRepo := db.NewUserRepository(database)
	roleRepo := db.NewRoleRepository(database)
	permissionRepo := db.NewPermissionRepository(database)
	taskRepo := db.NewTaskRepository(database)
	projectRepo := db.NewProjectRepository(database)
	wbsRepo := db.NewWBSRepository(database)
	stakeholderRepo := db.NewStakeholderRepository(database)
	holidayRepo := db.NewHolidayRepository(database)
	settingRepo := db.NewSettingRepository(database)
	categoryRepo := db.NewCategoryRepository(database)
	resourceRepo := db.NewResourceRepository(database)
	riskRepo := db.NewRiskRepository(database)
	issueRepo := db.NewIssueRepository(database)
	timesheetRepo := db.NewTimesheetRepository(database)
	projectMemberRepo := db.NewProjectMemberRepository(database)
	projectRoleRepo := db.NewProjectRoleRepository(database)
	snapshotRepo := db.NewSnapshotRepository(database)
	notificationRepo := db.NewNotificationRepository(database)

	authService, _ := service.NewAuthService(userRepo, roleRepo, "../certs/private.pem", "../certs/public.pem")
	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	wbsService := service.NewWBSService(wbsRepo)
	stakeholderService := service.NewStakeholderService(stakeholderRepo)
	holidayService := service.NewHolidayService(holidayRepo)
	settingService := service.NewSettingService(settingRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	portfolioService := service.NewPortfolioService(projectRepo)
	resourceService := service.NewResourceService(resourceRepo)
	riskService := service.NewRiskService(riskRepo)
	issueService := service.NewIssueService(issueRepo)
	snapshotService := service.NewSnapshotService(projectRepo, wbsRepo, snapshotRepo)
	reportService := service.NewReportService(projectRepo, wbsRepo)
	exportService := service.NewExportService(projectRepo, wbsRepo, snapshotRepo, riskRepo)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)
	timesheetService := service.NewTimesheetService(timesheetRepo, wbsRepo, notificationService)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	projectMemberService := service.NewProjectMemberService(projectMemberRepo)
	projectRoleService := service.NewProjectRoleService(projectRoleRepo)

	r := presentation.SetupRouter(authService, userRepo, taskService, userService, projectService, wbsService, stakeholderService, holidayService, settingService, categoryService, portfolioService, resourceService, riskService, issueService, reportService, exportService, timesheetService, notificationService, roleService, projectMemberService, projectRoleService, snapshotService)

	testServer = httptest.NewServer(r)
	defer testServer.Close()

	// 6. Run Tests
	code := m.Run()

	os.Exit(code)
}

func runMigrations(db *sql.DB) error {
	// Clean database
	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		return err
	}

	migrationDir := "../internal/infrastructure/db/migrations"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	var upFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			upFiles = append(upFiles, f.Name())
		}
	}
	sort.Strings(upFiles)

	for _, f := range upFiles {
		content, err := os.ReadFile(filepath.Join(migrationDir, f))
		if err != nil {
			return err
		}
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error in %s: %v", f, err)
		}
	}

	return nil
}
