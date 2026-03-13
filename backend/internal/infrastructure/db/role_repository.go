package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type roleRepo struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) repository.RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) List(ctx context.Context) ([]entity.Role, error) {
	query := `SELECT id, name, description, created_at FROM roles ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *roleRepo) GetByID(ctx context.Context, id uint) (*entity.Role, error) {
	role := &entity.Role{}
	query := `SELECT id, name, description, created_at FROM roles WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	role := &entity.Role{}
	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepo) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, role.Name, role.Description).Scan(&role.ID, &role.CreatedAt)
}

func (r *roleRepo) Update(ctx context.Context, role *entity.Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, role.Name, role.Description, role.ID)
	return err
}

func (r *roleRepo) Delete(ctx context.Context, id uint) error {
	// First delete role_permissions
	queryDeletePerms := `DELETE FROM role_permissions WHERE role_id = $1`
	_, err := r.db.ExecContext(ctx, queryDeletePerms, id)
	if err != nil {
		return err
	}

	query := `DELETE FROM roles WHERE id = $1`
	_, err = r.db.ExecContext(ctx, query, id)
	return err
}

func (r *roleRepo) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing permissions
	_, err = tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = $1`, roleID)
	if err != nil {
		return err
	}

	// Insert new ones
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`
	for _, pid := range permissionIDs {
		_, err = tx.ExecContext(ctx, query, roleID, pid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *roleRepo) GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]entity.Permission, error) {
	query := `SELECT p.id, p.name, p.description 
			  FROM permissions p 
			  JOIN role_permissions rp ON p.id = rp.permission_id 
			  WHERE rp.role_id = $1`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []entity.Permission
	for rows.Next() {
		var p entity.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}
