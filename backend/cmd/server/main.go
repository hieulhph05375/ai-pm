package main

import (
	"log"
	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/infrastructure/config"
	"project-mgmt/backend/internal/infrastructure/db"
	"project-mgmt/backend/internal/presentation"
)

func main() {
	cfg := config.Load()
	database, err := db.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

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

	authService, err := service.NewAuthService(userRepo, roleRepo, "certs/private.pem", "certs/public.pem")
	if err != nil {
		log.Fatalf("Error initializing auth service: %v", err)
	}

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
	snapshotRepo := db.NewSnapshotRepository(database)
	snapshotService := service.NewSnapshotService(projectRepo, wbsRepo, snapshotRepo)
	reportService := service.NewReportService(projectRepo, wbsRepo)
	exportService := service.NewExportService(projectRepo, wbsRepo, snapshotRepo, riskRepo)
	notificationRepo := db.NewNotificationRepository(database)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)
	timesheetService := service.NewTimesheetService(timesheetRepo, wbsRepo, notificationService)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	projectMemberService := service.NewProjectMemberService(projectMemberRepo)
	projectRoleService := service.NewProjectRoleService(projectRoleRepo)

	// Start background cron jobs
	service.StartCronJobs(notificationService, userRepo)

	r := presentation.SetupRouter(authService, userRepo, taskService, userService, projectService, wbsService, stakeholderService, holidayService, settingService, categoryService, portfolioService, resourceService, riskService, issueService, reportService, exportService, timesheetService, notificationService, roleService, projectMemberService, projectRoleService, snapshotService)
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
