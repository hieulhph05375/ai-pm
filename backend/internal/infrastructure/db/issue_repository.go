package db

import (
	"context"
	"database/sql"
	"project-mgmt/backend/internal/domain/entity"
)

type issueRepo struct {
	db *sql.DB
}

func NewIssueRepository(db *sql.DB) *issueRepo {
	return &issueRepo{db: db}
}

func nullStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	s := ns.String
	return &s
}

func (r *issueRepo) Create(ctx context.Context, i *entity.Issue) error {
	query := `
		INSERT INTO issues (project_id, type_id, title, description, status_id, priority_id, assignee_id, reporter_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		i.ProjectID, i.TypeID, i.Title, i.Description, i.StatusID, i.PriorityID, i.AssigneeID, i.ReporterID,
	).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
}

func (r *issueRepo) ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Issue, error) {
	query := `
		SELECT 
			i.id, i.project_id, i.title, i.description, i.assignee_id, i.reporter_id, i.created_at, i.updated_at,
			i.type_id, t.name as type_name, t.color as type_color, t.icon as type_icon,
			i.priority_id, p.name as priority_name, p.color as priority_color, p.icon as priority_icon,
			i.status_id, s.name as status_name, s.color as status_color, s.icon as status_icon
		FROM issues i
		LEFT JOIN categories t ON i.type_id = t.id
		LEFT JOIN categories p ON i.priority_id = p.id
		LEFT JOIN categories s ON i.status_id = s.id
		WHERE i.project_id = $1
		ORDER BY
			CASE p.name WHEN 'Critical' THEN 1 WHEN 'High' THEN 2 WHEN 'Medium' THEN 3 ELSE 4 END,
			i.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	issues := []entity.Issue{}
	for rows.Next() {
		var i entity.Issue
		var typeID, priorityID, statusID sql.NullInt64
		var typeName, typeColor, typeIcon sql.NullString
		var priorityName, priorityColor, priorityIcon sql.NullString
		var statusName, statusColor, statusIcon sql.NullString

		if err := rows.Scan(
			&i.ID, &i.ProjectID, &i.Title, &i.Description, &i.AssigneeID, &i.ReporterID, &i.CreatedAt, &i.UpdatedAt,
			&typeID, &typeName, &typeColor, &typeIcon,
			&priorityID, &priorityName, &priorityColor, &priorityIcon,
			&statusID, &statusName, &statusColor, &statusIcon,
		); err != nil {
			return nil, err
		}

		if typeID.Valid {
			i.TypeID = uint(typeID.Int64)
			i.Type = &entity.Category{
				ID:    i.TypeID,
				Name:  typeName.String,
				Color: nullStringPtr(typeColor),
				Icon:  nullStringPtr(typeIcon),
			}
		}
		if priorityID.Valid {
			i.PriorityID = uint(priorityID.Int64)
			i.Priority = &entity.Category{
				ID:    i.PriorityID,
				Name:  priorityName.String,
				Color: nullStringPtr(priorityColor),
				Icon:  nullStringPtr(priorityIcon),
			}
		}
		if statusID.Valid {
			i.StatusID = uint(statusID.Int64)
			i.Status = &entity.Category{
				ID:    i.StatusID,
				Name:  statusName.String,
				Color: nullStringPtr(statusColor),
				Icon:  nullStringPtr(statusIcon),
			}
		}

		issues = append(issues, i)
	}
	return issues, nil
}

func (r *issueRepo) CountByProject(ctx context.Context, projectID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM issues WHERE project_id = $1", projectID).Scan(&count)
	return count, err
}

func (r *issueRepo) Update(ctx context.Context, i *entity.Issue) error {
	query := `
		UPDATE issues
		SET type_id=$1, title=$2, description=$3, status_id=$4, priority_id=$5, assignee_id=$6, updated_at=NOW()
		WHERE id=$7`
	_, err := r.db.ExecContext(ctx, query,
		i.TypeID, i.Title, i.Description, i.StatusID, i.PriorityID, i.AssigneeID, i.ID,
	)
	return err
}

func (r *issueRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM issues WHERE id=$1", id)
	return err
}
