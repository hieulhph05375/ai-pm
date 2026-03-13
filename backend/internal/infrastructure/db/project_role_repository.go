package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type ProjectRoleRepository interface {
	Create(ctx context.Context, role *entity.ProjectRole) error
	GetByID(ctx context.Context, id uint) (*entity.ProjectRole, error)
	GetByProject(ctx context.Context, projectID uint) ([]entity.ProjectRole, error)
	Update(ctx context.Context, role *entity.ProjectRole) error
	Delete(ctx context.Context, id uint) error

	GetPermissions(ctx context.Context) ([]entity.ProjectPermission, error)
	SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]entity.ProjectPermission, error)
}

type projectRoleRepository struct {
	db *sql.DB
}

func NewProjectRoleRepository(db *sql.DB) ProjectRoleRepository {
	return &projectRoleRepository{db: db}
}

func (r *projectRoleRepository) Create(ctx context.Context, role *entity.ProjectRole) error {
	query := `INSERT INTO project_roles (project_id, name, description, color, is_default) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, role.ProjectID, role.Name, role.Description, role.Color, role.IsDefault).
		Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)
}

func (r *projectRoleRepository) GetByID(ctx context.Context, id uint) (*entity.ProjectRole, error) {
	role := &entity.ProjectRole{}
	query := `SELECT id, project_id, name, description, color, is_default, created_at, updated_at FROM project_roles WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&role.ID, &role.ProjectID, &role.Name, &role.Description, &role.Color, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return role, err
}

func (r *projectRoleRepository) GetByProject(ctx context.Context, projectID uint) ([]entity.ProjectRole, error) {
	query := `SELECT id, project_id, name, description, color, is_default, created_at, updated_at 
	          FROM project_roles WHERE project_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []entity.ProjectRole
	for rows.Next() {
		var role entity.ProjectRole
		if err := rows.Scan(&role.ID, &role.ProjectID, &role.Name, &role.Description, &role.Color, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *projectRoleRepository) Update(ctx context.Context, role *entity.ProjectRole) error {
	query := `UPDATE project_roles SET name = $1, description = $2, color = $3, is_default = $4, updated_at = CURRENT_TIMESTAMP 
	          WHERE id = $5`
	_, err := r.db.ExecContext(ctx, query, role.Name, role.Description, role.Color, role.IsDefault, role.ID)
	return err
}

func (r *projectRoleRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM project_roles WHERE id = $1", id)
	return err
}

func (r *projectRoleRepository) GetPermissions(ctx context.Context) ([]entity.ProjectPermission, error) {
	query := `SELECT id, name, description, module, created_at FROM project_permissions ORDER BY module, name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []entity.ProjectPermission
	for rows.Next() {
		var p entity.ProjectPermission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Module, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *projectRoleRepository) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM project_role_permissions WHERE project_role_id = $1", roleID); err != nil {
		return err
	}

	for _, pid := range permissionIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO project_role_permissions (project_role_id, project_permission_id) VALUES ($1, $2)", roleID, pid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *projectRoleRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]entity.ProjectPermission, error) {
	query := `SELECT p.id, p.name, p.description, p.module, p.created_at 
	          FROM project_permissions p
	          JOIN project_role_permissions prp ON p.id = prp.project_permission_id
	          WHERE prp.project_role_id = $1`
	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []entity.ProjectPermission
	for rows.Next() {
		var p entity.ProjectPermission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Module, &p.CreatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}
