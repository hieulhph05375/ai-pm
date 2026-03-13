package db

import (
	"context"
	"database/sql"
	"time"

	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type timesheetRepo struct {
	db *sql.DB
}

func NewTimesheetRepository(db *sql.DB) repository.TimesheetRepository {
	return &timesheetRepo{db: db}
}

func (r *timesheetRepo) Create(ctx context.Context, t *entity.Timesheet) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	if t.Status == "" {
		t.Status = entity.TimesheetStatusDraft
	}

	query := `INSERT INTO timesheets (user_id, project_id, node_id, task_id, work_date, hours, description, status, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		t.UserID, t.ProjectID, t.NodeID, t.TaskID, t.WorkDate, t.Hours, t.Description, t.Status, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
}

func (r *timesheetRepo) GetByID(ctx context.Context, id int) (*entity.Timesheet, error) {
	t := &entity.Timesheet{}
	query := `
		SELECT t.id, t.user_id, t.project_id, t.node_id, t.task_id, t.work_date, t.hours, t.description, t.status, t.created_at, t.updated_at,
			   u.full_name as user_name,
			   COALESCE(p.project_name, '') as project_name,
			   COALESCE(w.title, '') as node_title,
			   COALESCE(tk.title, '') as task_title
		FROM timesheets t
		JOIN users u ON t.user_id = u.id
		LEFT JOIN projects p ON t.project_id = p.id
		LEFT JOIN wbs_nodes w ON t.node_id = w.id
		LEFT JOIN tasks tk ON t.task_id = tk.id
		WHERE t.id = $1`

	var projName, nodeTitle, taskTitle sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.UserID, &t.ProjectID, &t.NodeID, &t.TaskID, &t.WorkDate, &t.Hours, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt,
		&t.UserName, &projName, &nodeTitle, &taskTitle,
	)
	if err != nil {
		return nil, err
	}
	t.ProjectName = projName.String
	t.NodeTitle = nodeTitle.String
	t.TaskTitle = taskTitle.String

	return t, nil
}

func (r *timesheetRepo) Update(ctx context.Context, t *entity.Timesheet) error {
	t.UpdatedAt = time.Now()
	query := `UPDATE timesheets SET
				project_id = $1, node_id = $2, task_id = $3, work_date = $4, hours = $5, description = $6, status = $7, updated_at = $8
			  WHERE id = $9`

	_, err := r.db.ExecContext(ctx, query,
		t.ProjectID, t.NodeID, t.TaskID, t.WorkDate, t.Hours, t.Description, t.Status, t.UpdatedAt, t.ID)
	return err
}

func (r *timesheetRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM timesheets WHERE id = $1`, id)
	return err
}

func (r *timesheetRepo) ListByUser(ctx context.Context, userID int, limit, offset int) ([]entity.Timesheet, int, error) {
	countQuery := `SELECT COUNT(*) FROM timesheets WHERE user_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT t.id, t.user_id, t.project_id, t.node_id, t.task_id, t.work_date, t.hours, t.description, t.status, t.created_at, t.updated_at,
			   u.full_name as user_name,
			   COALESCE(p.project_name, '') as project_name,
			   COALESCE(w.title, '') as node_title,
			   COALESCE(tk.title, '') as task_title
		FROM timesheets t
		JOIN users u ON t.user_id = u.id
		LEFT JOIN projects p ON t.project_id = p.id
		LEFT JOIN wbs_nodes w ON t.node_id = w.id
		LEFT JOIN tasks tk ON t.task_id = tk.id
		WHERE t.user_id = $1
		ORDER BY t.work_date DESC, t.created_at DESC
		LIMIT $2 OFFSET $3`
	items, _, err := r.list(ctx, query, userID, limit, offset)
	return items, total, err
}

func (r *timesheetRepo) ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Timesheet, int, error) {
	countQuery := `SELECT COUNT(*) FROM timesheets WHERE project_id = $1`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT t.id, t.user_id, t.project_id, t.node_id, t.task_id, t.work_date, t.hours, t.description, t.status, t.created_at, t.updated_at,
			   u.full_name as user_name,
			   COALESCE(p.project_name, '') as project_name,
			   COALESCE(w.title, '') as node_title,
			   COALESCE(tk.title, '') as task_title
		FROM timesheets t
		JOIN users u ON t.user_id = u.id
		LEFT JOIN projects p ON t.project_id = p.id
		LEFT JOIN wbs_nodes w ON t.node_id = w.id
		LEFT JOIN tasks tk ON t.task_id = tk.id
		WHERE t.project_id = $1
		ORDER BY t.work_date DESC, t.created_at DESC
		LIMIT $2 OFFSET $3`
	items, _, err := r.list(ctx, query, projectID, limit, offset)
	return items, total, err
}

func (r *timesheetRepo) list(ctx context.Context, query string, arg1 int, limit, offset int) ([]entity.Timesheet, int, error) {
	rows, err := r.db.QueryContext(ctx, query, arg1, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var timesheets []entity.Timesheet
	for rows.Next() {
		var t entity.Timesheet
		var projName, nodeTitle, taskTitle sql.NullString
		err := rows.Scan(
			&t.ID, &t.UserID, &t.ProjectID, &t.NodeID, &t.TaskID, &t.WorkDate, &t.Hours, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt,
			&t.UserName, &projName, &nodeTitle, &taskTitle,
		)
		if err != nil {
			return nil, 0, err
		}
		t.ProjectName = projName.String
		t.NodeTitle = nodeTitle.String
		t.TaskTitle = taskTitle.String
		timesheets = append(timesheets, t)
	}

	return timesheets, 0, nil
}
