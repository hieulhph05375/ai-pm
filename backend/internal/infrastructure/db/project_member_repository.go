package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type projectMemberRepo struct {
	db *sql.DB
}

func NewProjectMemberRepository(db *sql.DB) *projectMemberRepo {
	return &projectMemberRepo{db: db}
}

func (r *projectMemberRepo) scanMember(rows interface {
	Scan(dest ...interface{}) error
}, m *entity.ProjectMember) error {
	return rows.Scan(
		&m.ID, &m.ProjectID, &m.UserID, &m.ProjectRoleID, &m.JoinedAt, &m.CreatedAt,
	)
}

func (r *projectMemberRepo) AddMember(ctx context.Context, member *entity.ProjectMember) error {
	query := `
		INSERT INTO project_members (project_id, user_id, project_role_id)
		VALUES ($1, $2, $3)
		RETURNING id, joined_at, created_at`
	return r.db.QueryRowContext(ctx, query, member.ProjectID, member.UserID, member.ProjectRoleID).
		Scan(&member.ID, &member.JoinedAt, &member.CreatedAt)
}

func (r *projectMemberRepo) RemoveMember(ctx context.Context, projectID, userID int) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, projectID, userID)
	return err
}

func (r *projectMemberRepo) UpdateMemberRole(ctx context.Context, projectID, userID int, roleID int) error {
	query := `UPDATE project_members SET project_role_id = $1 WHERE project_id = $2 AND user_id = $3`
	_, err := r.db.ExecContext(ctx, query, roleID, projectID, userID)
	return err
}

func (r *projectMemberRepo) GetMembersByProject(ctx context.Context, projectID int, page, limit int) ([]entity.ProjectMember, int, error) {
	// Count total members
	countQuery := `SELECT COUNT(*) FROM project_members WHERE project_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, projectID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated members
	offset := (page - 1) * limit
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.project_role_id, pm.joined_at, pm.created_at, 
		       u.full_name, u.email, 
		       pr.name, pr.color
		FROM project_members pm
		JOIN users u ON pm.user_id = u.id
		JOIN project_roles pr ON pm.project_role_id = pr.id
		WHERE pm.project_id = $1
		ORDER BY pm.joined_at DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	members := []entity.ProjectMember{}
	for rows.Next() {
		var m entity.ProjectMember
		var user entity.User
		var role entity.ProjectRole
		err := rows.Scan(
			&m.ID, &m.ProjectID, &m.UserID, &m.ProjectRoleID, &m.JoinedAt, &m.CreatedAt,
			&user.FullName, &user.Email,
			&role.Name, &role.Color,
		)
		if err != nil {
			return nil, 0, err
		}
		user.ID = uint(m.UserID)
		m.User = &user
		role.ID = uint(m.ProjectRoleID)
		m.Role = &role
		members = append(members, m)
	}
	return members, total, nil
}

func (r *projectMemberRepo) GetProjectsByUser(ctx context.Context, userID int) ([]entity.ProjectMember, error) {
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.project_role_id, pm.joined_at, pm.created_at, 
		       p.project_name, p.project_id as biz_id,
		       pr.name, pr.color
		FROM project_members pm
		JOIN projects p ON pm.project_id = p.id
		JOIN project_roles pr ON pm.project_role_id = pr.id
		WHERE pm.user_id = $1
		ORDER BY pm.joined_at DESC`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []entity.ProjectMember{}
	for rows.Next() {
		var m entity.ProjectMember
		var p entity.Project
		var role entity.ProjectRole
		err := rows.Scan(
			&m.ID, &m.ProjectID, &m.UserID, &m.ProjectRoleID, &m.JoinedAt, &m.CreatedAt,
			&p.ProjectName, &p.ProjectID,
			&role.Name, &role.Color,
		)
		if err != nil {
			return nil, err
		}
		p.ID = m.ProjectID
		m.Project = &p
		role.ID = uint(m.ProjectRoleID)
		m.Role = &role
		memberships = append(memberships, m)
	}
	return memberships, nil
}

func (r *projectMemberRepo) IsMember(ctx context.Context, projectID, userID int) (bool, int, error) {
	query := `SELECT project_role_id FROM project_members WHERE project_id = $1 AND user_id = $2`
	var roleID int
	err := r.db.QueryRowContext(ctx, query, projectID, userID).Scan(&roleID)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, roleID, nil
}
