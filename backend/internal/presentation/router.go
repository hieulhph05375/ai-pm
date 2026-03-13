package presentation

import (
	"os"

	"project-mgmt/backend/internal/application/service"
	"project-mgmt/backend/internal/domain/repository"
	"project-mgmt/backend/internal/presentation/handler"
	"project-mgmt/backend/internal/presentation/middleware"

	"github.com/gin-contrib/gzip"
	gingpprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authService service.AuthService, userRepo repository.UserRepository, taskService service.TaskService, userService service.UserService, projectService service.ProjectService, wbsService service.WBSService, stakeholderService service.StakeholderService, holidayService service.HolidayService, settingService service.SettingService, categoryService service.CategoryService, portfolioService service.PortfolioService, resourceService service.ResourceService, riskService service.RiskService, issueService service.IssueService, reportService service.ReportService, exportService service.ExportService, timesheetService service.TimesheetService, notificationService service.NotificationService, roleService service.RoleService, projectMemberService service.ProjectMemberService, projectRoleService service.ProjectRoleService, snapshotService service.SnapshotService) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.JSONLogMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Enable pprof in development mode only — never expose in production
	if os.Getenv("APP_ENV") != "production" {
		gingpprof.Register(r)
	}

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	taskHandler := handler.NewTaskHandler(taskService)
	projectHandler := handler.NewProjectHandler(projectService)
	wbsHandler := handler.NewWBSHandler(wbsService)
	stakeholderHandler := handler.NewStakeholderHandler(stakeholderService)
	holidayHandler := handler.NewHolidayHandler(holidayService)
	settingHandler := handler.NewSettingHandler(settingService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	portfolioHandler := handler.NewPortfolioHandler(portfolioService)
	resourceHandler := handler.NewResourceHandler(resourceService)
	riskHandler := handler.NewRiskHandler(riskService)
	issueHandler := handler.NewIssueHandler(issueService)
	reportHandler := handler.NewReportHandler(reportService, exportService)
	timesheetHandler := handler.NewTimesheetHandler(timesheetService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	roleHandler := handler.NewRoleHandler(roleService)
	projectMemberHandler := handler.NewProjectMemberHandler(projectMemberService)
	projectRoleHandler := handler.NewProjectRoleHandler(projectRoleService)
	reportingHandler := handler.NewReportingHandler(snapshotService)

	authMW := middleware.NewAuthorizer(authService.GetPublicKey())

	v1 := r.Group("/api/v1")
	{
		// Health Check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "UP",
			})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.Refresh)
		}

		users := v1.Group("/users")
		users.Use(authMW.Authenticate())
		{
			users.GET("/me", userHandler.GetMe)
			users.GET("", middleware.AuthorizePermission("user:read"), userHandler.List)
			users.POST("", middleware.AuthorizePermission("user:create"), userHandler.Create)
			users.PUT("/:id", middleware.AuthorizePermission("user:update"), userHandler.Update)
			users.PUT("/:id/reset-password", middleware.AuthorizePermission("user:update"), userHandler.ResetPassword)
			users.PUT("/:id/toggle-status", middleware.AuthorizePermission("user:update"), userHandler.ToggleStatus)
		}

		tasks := v1.Group("/tasks")
		tasks.Use(authMW.Authenticate())
		{
			tasks.GET("", middleware.AuthorizePermission("task:read"), taskHandler.List)
			tasks.POST("", middleware.AuthorizePermission("task:create"), taskHandler.Create)
			tasks.GET("/:id", middleware.AuthorizePermission("task:read"), taskHandler.GetByID)
			tasks.PUT("/:id", middleware.AuthorizePermission("task:update"), taskHandler.Update)
			tasks.DELETE("/:id", middleware.AuthorizePermission("task:delete"), taskHandler.Delete)
			tasks.GET("/:id/activities", middleware.AuthorizePermission("task:read"), taskHandler.ListActivities)
			tasks.POST("/:id/comments", middleware.AuthorizePermission("task:update"), taskHandler.AddComment)
		}

		// Stakeholder Routes (Global)
		stakeholders := v1.Group("/stakeholders")
		stakeholders.Use(authMW.Authenticate())
		{
			stakeholders.GET("", middleware.AuthorizePermission("stakeholder:read"), stakeholderHandler.List)
			stakeholders.POST("", middleware.AuthorizePermission("stakeholder:create"), stakeholderHandler.Create)
			stakeholders.GET("/:id", middleware.AuthorizePermission("stakeholder:read"), stakeholderHandler.GetByID)
			stakeholders.PUT("/:id", middleware.AuthorizePermission("stakeholder:update"), stakeholderHandler.Update)
			stakeholders.DELETE("/:id", middleware.AuthorizePermission("stakeholder:delete"), stakeholderHandler.Delete)
		}

		// Settings Routes (Global)
		settings := v1.Group("/settings")
		settings.Use(authMW.Authenticate())
		{
			settings.GET("", middleware.AuthorizePermission("setting:read"), settingHandler.GetAll)
			settings.PUT("/:key", middleware.AuthorizePermission("setting:update"), settingHandler.Update)
		}

		// Holiday Routes (Global)
		holidays := v1.Group("/holidays")
		holidays.Use(authMW.Authenticate())
		{
			holidays.GET("", middleware.AuthorizePermission("holiday:read"), holidayHandler.List)
			holidays.POST("", middleware.AuthorizePermission("holiday:create"), holidayHandler.Create)
			holidays.GET("/:id", middleware.AuthorizePermission("holiday:read"), holidayHandler.GetByID)
			holidays.PUT("/:id", middleware.AuthorizePermission("holiday:update"), holidayHandler.Update)
			holidays.DELETE("/:id", middleware.AuthorizePermission("holiday:delete"), holidayHandler.Delete)
		}

		// Category Management Routes (Global)
		categoryTypes := v1.Group("/category-types")
		categoryTypes.Use(authMW.Authenticate())
		{
			categoryTypes.GET("", middleware.AuthorizePermission("category:read"), categoryHandler.ListTypes)
			categoryTypes.POST("", middleware.AuthorizePermission("category:create"), categoryHandler.CreateType)
			categoryTypes.PUT("/:id", middleware.AuthorizePermission("category:update"), categoryHandler.UpdateType)
			categoryTypes.DELETE("/:id", middleware.AuthorizePermission("category:delete"), categoryHandler.DeleteType)
		}

		categories := v1.Group("/categories")
		categories.Use(authMW.Authenticate())
		{
			categories.GET("", middleware.AuthorizePermission("category:read"), categoryHandler.ListCategories)
			categories.POST("", middleware.AuthorizePermission("category:create"), categoryHandler.CreateCategory)
			categories.GET("/:id", middleware.AuthorizePermission("category:read"), categoryHandler.GetCategoryByID)
			categories.PUT("/:id", middleware.AuthorizePermission("category:update"), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.AuthorizePermission("category:delete"), categoryHandler.DeleteCategory)
		}

		projects := v1.Group("/projects")
		projects.Use(authMW.Authenticate())
		{
			projects.GET("", middleware.AuthorizePermission("project:read"), projectHandler.List)
			projects.GET("/export", middleware.AuthorizePermission("project:export"), reportHandler.ExportProjectList)
			projects.POST("", middleware.AuthorizePermission("project:create"), projectHandler.Create)
			projects.GET("/:id", authMW.AuthorizeProjectAccess(projectMemberService), middleware.AuthorizePermission("project:read"), projectHandler.GetByID)
			projects.PUT("/:id", authMW.AuthorizeProjectAccess(projectMemberService), middleware.AuthorizePermission("project:update"), projectHandler.Update)
			projects.DELETE("/:id", middleware.AuthorizePermission("project:delete"), projectHandler.Delete)

			// Project Stakeholder Mapping
			projectSpecific := projects.Group("/:id")
			projectSpecific.Use(authMW.AuthorizeProjectAccess(projectMemberService))
			{
				projectSpecific.GET("/stakeholders", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:read"), stakeholderHandler.ListByProject)
				projectSpecific.POST("/stakeholders", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:update"), stakeholderHandler.AssignToProject)
				projectSpecific.DELETE("/stakeholders/:sid", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:update"), stakeholderHandler.UnassignFromProject)

				// Project Membership Routes
				projectSpecific.GET("/members", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:team:view"), projectMemberHandler.GetMembers)
				projectSpecific.POST("/members", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:team:create"), projectMemberHandler.AddMember)
				projectSpecific.PUT("/members/:userId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:team:update"), projectMemberHandler.UpdateMemberRole)
				projectSpecific.DELETE("/members/:userId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:team:delete"), projectMemberHandler.RemoveMember)

				// Project Role & Permission Routes
				projectSpecific.GET("/permissions", projectRoleHandler.GetPermissions)
				projectSpecific.GET("/roles", projectRoleHandler.GetRolesByProject)
				projectSpecific.POST("/roles", middleware.AuthorizePermission("project:roles:create"), projectRoleHandler.CreateRole)
				projectSpecific.PUT("/roles/:roleId", middleware.AuthorizePermission("project:roles:update"), projectRoleHandler.UpdateRole)
				projectSpecific.DELETE("/roles/:roleId", middleware.AuthorizePermission("project:roles:delete"), projectRoleHandler.DeleteRole)
				projectSpecific.GET("/roles/:roleId/permissions", projectRoleHandler.GetRolePermissions)
				projectSpecific.PUT("/roles/:roleId/permissions", middleware.AuthorizePermission("project:roles:update"), projectRoleHandler.SetPermissions)

				// WBS Routes
				projectSpecific.GET("/wbs", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:view"), wbsHandler.ListTree)
				projectSpecific.POST("/wbs", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:create"), wbsHandler.CreateNode)
				projectSpecific.GET("/wbs/:nodeId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:view"), wbsHandler.GetNode)
				projectSpecific.PUT("/wbs/:nodeId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:update"), wbsHandler.UpdateNode)
				projectSpecific.DELETE("/wbs/:nodeId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:delete"), wbsHandler.DeleteNode)

				projectSpecific.GET("/wbs/dependencies", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:view"), wbsHandler.ListDependencies)
				projectSpecific.POST("/wbs/dependencies", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:create"), wbsHandler.CreateDependency)
				projectSpecific.DELETE("/wbs/dependencies/:depId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:delete"), wbsHandler.DeleteDependency)

				// WBS Comments
				projectSpecific.GET("/wbs/:nodeId/comments", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:view"), wbsHandler.ListComments)
				projectSpecific.POST("/wbs/:nodeId/comments", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:update"), wbsHandler.AddComment)
				projectSpecific.PUT("/wbs/:nodeId/comments/:commentId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:update"), wbsHandler.UpdateComment)
				projectSpecific.DELETE("/wbs/:nodeId/comments/:commentId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:delete"), wbsHandler.DeleteComment)

				// WBS Baselines
				projectSpecific.GET("/wbs-baselines", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:read"), wbsHandler.ListBaselines)
				projectSpecific.POST("/wbs-baselines", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:update"), wbsHandler.CreateBaseline)
				projectSpecific.GET("/wbs-baselines/:baselineId/nodes", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:read"), wbsHandler.GetBaselineNodes)

				// Risk Register
				projectSpecific.GET("/risks", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:risk:view"), riskHandler.List)
				projectSpecific.POST("/risks", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:risk:create"), riskHandler.Create)
				projectSpecific.PUT("/risks/:riskId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:risk:update"), riskHandler.Update)
				projectSpecific.DELETE("/risks/:riskId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:risk:delete"), riskHandler.Delete)

				// Issue & Bug Tracker
				projectSpecific.GET("/issues", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:issue:view"), issueHandler.List)
				projectSpecific.POST("/issues", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:issue:create"), issueHandler.Create)
				projectSpecific.PUT("/issues/:issueId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:issue:update"), issueHandler.Update)
				projectSpecific.DELETE("/issues/:issueId", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:issue:delete"), issueHandler.Delete)

				// PMI & Reports
				projectSpecific.GET("/pmi-stats", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:read"), reportHandler.GetPMIStats)
				projectSpecific.GET("/export/wbs", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:wbs:view"), reportHandler.ExportWBS)
				projectSpecific.GET("/export/summary", authMW.AuthorizeProjectPermission(projectMemberService, projectRoleService, "project:read"), reportHandler.ExportSummary)
			}
		}

		// Portfolio Routes (PMO View)
		portfolio := v1.Group("/portfolio")
		portfolio.Use(authMW.Authenticate())
		{
			portfolio.GET("/overview", middleware.AuthorizePermission("portfolio:read"), portfolioHandler.GetOverview)
		}

		// Resource Allocation Routes (Workload Heatmap)
		resources := v1.Group("/resources")
		resources.Use(authMW.Authenticate())
		{
			resources.GET("/workload", middleware.AuthorizePermission("resource:read"), resourceHandler.GetWorkload)
		}

		// Reporting & Analytics
		reporting := v1.Group("/reporting")
		reporting.Use(authMW.Authenticate())
		{
			reporting.POST("/snapshots/capture", middleware.AuthorizePermission("portfolio:read"), reportingHandler.CaptureSnapshots)
			reporting.GET("/projects/:id/trends", authMW.AuthorizeProjectAccess(projectMemberService), middleware.AuthorizePermission("project:read"), reportingHandler.GetProjectTrends)
			reporting.GET("/projects/:id/milestone-trends", authMW.AuthorizeProjectAccess(projectMemberService), middleware.AuthorizePermission("project:read"), reportingHandler.GetMilestoneTrends)
		}

		// Timesheet Routes
		timesheets := v1.Group("/timesheets")
		timesheets.Use(authMW.Authenticate())
		{
			timesheets.GET("", middleware.AuthorizePermission("timesheet:read"), timesheetHandler.List)
			timesheets.POST("", middleware.AuthorizePermission("timesheet:create"), timesheetHandler.Create)
			timesheets.GET("/:id", middleware.AuthorizePermission("timesheet:read"), timesheetHandler.GetByID)
			timesheets.PUT("/:id", middleware.AuthorizePermission("timesheet:update"), timesheetHandler.Update)
			timesheets.DELETE("/:id", middleware.AuthorizePermission("timesheet:delete"), timesheetHandler.Delete)
		}

		// Notification Routes
		notifications := v1.Group("/notifications")
		notifications.Use(authMW.Authenticate())
		{
			notifications.GET("", notificationHandler.List)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllRead)
		}

		// RBAC Management Routes (Admin only fallback, now enforced by permissions)
		roles := v1.Group("/roles")
		roles.Use(authMW.Authenticate())
		{
			roles.GET("", middleware.AuthorizePermission("role:read"), roleHandler.ListRoles)
			roles.GET("/:id", middleware.AuthorizePermission("role:read"), roleHandler.GetRoleWithPermissions)
			roles.POST("", middleware.AuthorizePermission("role:create"), roleHandler.CreateRole)
			roles.PUT("/:id", middleware.AuthorizePermission("role:update"), roleHandler.UpdateRole)
			roles.DELETE("/:id", middleware.AuthorizePermission("role:delete"), roleHandler.DeleteRole)
			roles.PUT("/:id/permissions", middleware.AuthorizePermission("role:update"), roleHandler.AssignPermissions)
		}

		permissions := v1.Group("/permissions")
		permissions.Use(authMW.Authenticate())
		{
			permissions.GET("", middleware.AuthorizePermission("role:read"), roleHandler.ListPermissions)
		}
	}
	return r
}
