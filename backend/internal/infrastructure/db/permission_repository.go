package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type permissionRepo struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &permissionRepo{db: db}
}

func (r *permissionRepo) List(ctx context.Context) ([]entity.Permission, error) {
	query := `SELECT id, name, description FROM permissions ORDER BY name ASC`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *permissionRepo) GetByID(ctx context.Context, id uint) (*entity.Permission, error) {
	p := &entity.Permission{}
	query := `SELECT id, name, description FROM permissions WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *permissionRepo) GetByName(ctx context.Context, name string) (*entity.Permission, error) {
	p := &entity.Permission{}
	query := `SELECT id, name, description FROM permissions WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&p.ID, &p.Name, &p.Description)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *permissionRepo) Create(ctx context.Context, p *entity.Permission) error {
	query := `INSERT INTO permissions (name, description) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.Name, p.Description).Scan(&p.ID)
}
