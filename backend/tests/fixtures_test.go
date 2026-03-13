package tests

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedAll(db *sql.DB) error {
	// 1. Clean tables to ensure fresh start
	tables := []string{
		"project_snapshots", "milestone_snapshots", "wbs_nodes", "project_members",
		"project_role_permissions", "project_roles", "projects", "role_permissions",
		"users", "roles", "permissions",
	}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %v", table, err)
		}
	}

	// 2. Roles (Standardized for RBAC tests)
	allRoles := []string{"admin", "pmo", "projectmanager", "teamlead", "member", "viewer"}
	roleIDMap := make(map[string]int)
	for _, r := range allRoles {
		var id int
		err := db.QueryRow("INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id", r, "Role "+r).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %v", r, err)
		}
		roleIDMap[r] = id
	}

	// 3. Create Users for each role
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	userIDMap := make(map[string]int)

	for roleName, roleID := range roleIDMap {
		email := roleName + "@example.com"
		var userID int
		err := db.QueryRow(`
			INSERT INTO users (email, hashed_password, full_name, role_id, is_active, created_at, updated_at, is_admin)
			VALUES ($1, $2, $3, $4, true, NOW(), NOW(), $5)
			RETURNING id`,
			email, string(passwordHash), roleName+" User", roleID, roleName == "admin").Scan(&userID)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", email, err)
		}
		userIDMap[roleName] = userID
	}

	// 4. Seed Permissions
	permissions := []string{
		"user:read", "user:create", "user:update", "user:delete",
		"role:read", "role:create", "role:update", "role:delete",
		"project:read", "project:create", "project:update", "project:delete",
		"risk:read", "risk:create", "risk:update", "risk:delete",
		"issue:read", "issue:create", "issue:update", "issue:delete",
		"timesheet:read", "timesheet:create", "timesheet:update", "timesheet:delete",
		"holiday:read", "holiday:create", "holiday:update", "holiday:delete",
		"setting:read", "setting:update",
		"reporting:read", "reporting:snapshot",
		"dashboard:read",
		"portfolio:read", "resource:read",
	}

	permMap := make(map[string]int)
	for _, p := range permissions {
		var id int
		err := db.QueryRow("INSERT INTO permissions (name, description) VALUES ($1, $2) RETURNING id", p, "Permission for "+p).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to seed permission %s: %v", p, err)
		}
		permMap[p] = id
	}

	// 5. Link Roles to Permissions
	rolePermMapping := map[string][]string{
		"admin":          permissions, // Everything
		"pmo":            {"project:read", "portfolio:read", "resource:read", "reporting:read", "reporting:snapshot", "dashboard:read", "setting:read", "holiday:read"},
		"projectmanager": {"project:read", "project:create", "project:update", "risk:read", "issue:read", "timesheet:read", "holiday:read", "setting:read", "resource:read"},
		"teamlead":       {"project:read", "timesheet:read", "holiday:read"},
		"member":         {"project:read", "timesheet:create", "holiday:read"},
		"viewer":         {"project:read", "holiday:read"},
	}

	for roleName, perms := range rolePermMapping {
		roleID := roleIDMap[roleName]
		for _, p := range perms {
			_, err := db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)", roleID, permMap[p])
			if err != nil {
				return fmt.Errorf("failed to link role %s to permission %s: %v", roleName, p, err)
			}
		}
	}

	// 6. Create Sample Project
	var projectID int
	err := db.QueryRow(`
		INSERT INTO projects (project_id, project_name, description, project_status, planned_start_date, planned_end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id`,
		"PRJ-001", "Test Project Alpha", "A project for testing", "Running", time.Now(), time.Now().AddDate(0, 3, 0)).Scan(&projectID)
	if err != nil {
		return fmt.Errorf("failed to create sample project: %v", err)
	}

	// 7. Seed Project Roles (Specific to project Alpha)
	var projPMID int
	err = db.QueryRow(`
		INSERT INTO project_roles (project_id, name, description, color, is_default)
		VALUES ($1, 'Project Manager', 'Project manager for Alpha', '#FF0000', true)
		RETURNING id`, projectID).Scan(&projPMID)
	if err != nil {
		return fmt.Errorf("failed to seed project role: %v", err)
	}

	// 8. Project Member (Admin is PM in project Alpha)
	_, err = db.Exec(`
		INSERT INTO project_members (project_id, user_id, project_role_id)
		VALUES ($1, $2, $3)`, projectID, userIDMap["admin"], projPMID)
	if err != nil {
		return fmt.Errorf("failed to add project member: %v", err)
	}

	// 9. Seed WBS Nodes
	_, err = db.Exec(`
		INSERT INTO wbs_nodes (project_id, title, path, type, progress, created_at, updated_at)
		VALUES 
		($1, 'Phase 1', 'P1', 'Phase', 0, NOW(), NOW()),
		($1, 'Task 1.1', 'P1.T1', 'Task', 0, NOW(), NOW())`,
		projectID)
	if err != nil {
		return fmt.Errorf("failed to create WBS nodes: %v", err)
	}

	return nil
}
