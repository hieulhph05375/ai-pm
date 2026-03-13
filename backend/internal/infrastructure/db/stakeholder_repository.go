package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
	"strconv"
)

type stakeholderRepo struct {
	db *sql.DB
}

func NewStakeholderRepository(db *sql.DB) *stakeholderRepo {
	return &stakeholderRepo{db: db}
}

func (r *stakeholderRepo) scanStakeholder(rows interface {
	Scan(dest ...interface{}) error
}, s *entity.Stakeholder) error {
	var roleID sql.NullInt64
	var roleName, roleColor, roleIcon sql.NullString

	err := rows.Scan(
		&s.ID, &s.Name, &s.Role, &roleID, &s.Email, &s.Phone, &s.Organization, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
		&roleName, &roleColor, &roleIcon,
	)
	if err != nil {
		return err
	}

	if roleID.Valid {
		s.RoleID = uint(roleID.Int64)
		s.RoleCat = &entity.Category{
			ID:    s.RoleID,
			Name:  roleName.String,
			Color: nullStringPtr(roleColor),
			Icon:  nullStringPtr(roleIcon),
		}
	}
	return nil
}

func (r *stakeholderRepo) Create(ctx context.Context, s *entity.Stakeholder) error {
	var roleID sql.NullInt64
	if s.RoleID > 0 {
		roleID = sql.NullInt64{Int64: int64(s.RoleID), Valid: true}
	}
	query := `
		INSERT INTO stakeholders (name, role, role_id, email, phone, organization, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, s.Name, s.Role, roleID, s.Email, s.Phone, s.Organization, s.Notes).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *stakeholderRepo) GetByID(ctx context.Context, id int) (*entity.Stakeholder, error) {
	s := &entity.Stakeholder{}
	query := `
		SELECT s.id, s.name, s.role, s.role_id, s.email, s.phone, s.organization, s.notes, s.created_at, s.updated_at,
		       c.name as role_name, c.color as role_color, c.icon as role_icon
		FROM stakeholders s
		LEFT JOIN categories c ON s.role_id = c.id
		WHERE s.id = $1`
	if err := r.scanStakeholder(r.db.QueryRowContext(ctx, query, id), s); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *stakeholderRepo) Update(ctx context.Context, s *entity.Stakeholder) error {
	var roleID sql.NullInt64
	if s.RoleID > 0 {
		roleID = sql.NullInt64{Int64: int64(s.RoleID), Valid: true}
	}
	query := `
		UPDATE stakeholders 
		SET name=$1, role=$2, role_id=$3, email=$4, phone=$5, organization=$6, notes=$7, updated_at=CURRENT_TIMESTAMP
		WHERE id=$8`
	_, err := r.db.ExecContext(ctx, query, s.Name, s.Role, roleID, s.Email, s.Phone, s.Organization, s.Notes, s.ID)
	return err
}

func (r *stakeholderRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM stakeholders WHERE id=$1", id)
	return err
}

func (r *stakeholderRepo) List(ctx context.Context, search string) ([]*entity.Stakeholder, error) {
	query := `
		SELECT s.id, s.name, s.role, s.role_id, s.email, s.phone, s.organization, s.notes, s.created_at, s.updated_at,
		       c.name as role_name, c.color as role_color, c.icon as role_icon
		FROM stakeholders s
		LEFT JOIN categories c ON s.role_id = c.id`
	var args []interface{}
	if search != "" {
		query += " WHERE s.name ILIKE $1 OR s.organization ILIKE $1 OR s.role ILIKE $1"
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY s.name ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*entity.Stakeholder{}
	for rows.Next() {
		s := &entity.Stakeholder{}
		if err := r.scanStakeholder(rows, s); err != nil {
			return nil, err
		}
		results = append(results, s)
	}
	return results, nil
}

func (r *stakeholderRepo) ListWithPagination(ctx context.Context, search string, offset, limit int) ([]*entity.Stakeholder, int, error) {
	query := `
		SELECT s.id, s.name, s.role, s.role_id, s.email, s.phone, s.organization, s.notes, s.created_at, s.updated_at,
		       c.name as role_name, c.color as role_color, c.icon as role_icon
		FROM stakeholders s
		LEFT JOIN categories c ON s.role_id = c.id`
	countQuery := `SELECT COUNT(*) FROM stakeholders s`
	var args []interface{}
	if search != "" {
		filter := " WHERE s.name ILIKE $1 OR s.organization ILIKE $1 OR s.role ILIKE $1"
		query += filter
		countQuery += filter
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY s.name ASC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	results := []*entity.Stakeholder{}
	for rows.Next() {
		s := &entity.Stakeholder{}
		if err := r.scanStakeholder(rows, s); err != nil {
			return nil, 0, err
		}
		results = append(results, s)
	}
	return results, total, nil
}

func (r *stakeholderRepo) Count(ctx context.Context, search string) (int, error) {
	query := `SELECT COUNT(*) FROM stakeholders`
	var args []interface{}
	if search != "" {
		query += " WHERE name ILIKE $1 OR organization ILIKE $1 OR role ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *stakeholderRepo) AssignToProject(ctx context.Context, projectID int, stakeholderID int, projectRole string, roleID uint) error {
	query := `
		INSERT INTO project_stakeholders (project_id, stakeholder_id, project_role, role_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (project_id, stakeholder_id) DO UPDATE SET project_role = $3, role_id = $4`
	_, err := r.db.ExecContext(ctx, query, projectID, stakeholderID, projectRole, roleID)
	return err
}

func (r *stakeholderRepo) UnassignFromProject(ctx context.Context, projectID int, stakeholderID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM project_stakeholders WHERE project_id=$1 AND stakeholder_id=$2", projectID, stakeholderID)
	return err
}

func (r *stakeholderRepo) ListByProject(ctx context.Context, projectID int) ([]*entity.ProjectStakeholder, error) {
	query := `
		SELECT ps.project_id, ps.stakeholder_id, ps.project_role, ps.role_id as project_role_id, ps.created_at,
		       s.id, s.name, s.role, s.role_id, s.email, s.phone, s.organization, s.notes, s.created_at, s.updated_at,
			   c.name as role_name, c.color as role_color, c.icon as role_icon,
			   cp.name as proj_role_name, cp.color as proj_role_color, cp.icon as proj_role_icon
		FROM project_stakeholders ps
		JOIN stakeholders s ON ps.stakeholder_id = s.id
		LEFT JOIN categories c ON s.role_id = c.id
		LEFT JOIN categories cp ON ps.role_id = cp.id
		WHERE ps.project_id = $1
		ORDER BY s.name ASC`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*entity.ProjectStakeholder{}
	for rows.Next() {
		ps := &entity.ProjectStakeholder{Stakeholder: &entity.Stakeholder{}}
		var roleID, projRoleID sql.NullInt64
		var roleName, roleColor, roleIcon sql.NullString
		var projRoleName, projRoleColor, projRoleIcon sql.NullString
		err := rows.Scan(
			&ps.ProjectID, &ps.StakeholderID, &ps.ProjectRole, &projRoleID, &ps.CreatedAt,
			&ps.Stakeholder.ID, &ps.Stakeholder.Name, &ps.Stakeholder.Role, &roleID, &ps.Stakeholder.Email,
			&ps.Stakeholder.Phone, &ps.Stakeholder.Organization, &ps.Stakeholder.Notes,
			&ps.Stakeholder.CreatedAt, &ps.Stakeholder.UpdatedAt,
			&roleName, &roleColor, &roleIcon,
			&projRoleName, &projRoleColor, &projRoleIcon,
		)
		if err != nil {
			return nil, err
		}

		if roleID.Valid {
			ps.Stakeholder.RoleID = uint(roleID.Int64)
			ps.Stakeholder.RoleCat = &entity.Category{
				ID:    ps.Stakeholder.RoleID,
				Name:  roleName.String,
				Color: nullStringPtr(roleColor),
				Icon:  nullStringPtr(roleIcon),
			}
		}

		if projRoleID.Valid {
			rid := uint(projRoleID.Int64)
			ps.RoleID = &rid
			ps.RoleCat = &entity.Category{
				ID:    rid,
				Name:  projRoleName.String,
				Color: nullStringPtr(projRoleColor),
				Icon:  nullStringPtr(projRoleIcon),
			}
		}
		results = append(results, ps)
	}
	return results, nil
}
