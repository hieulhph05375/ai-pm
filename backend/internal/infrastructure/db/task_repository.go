package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type taskRepo struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(ctx context.Context, t *entity.Task) error {
	labelsJSON, err := json.Marshal(t.Labels)
	if err != nil {
		labelsJSON = []byte("[]")
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `INSERT INTO tasks (title, description, status, priority, status_id, priority_id, assignee_id, created_by, start_date, due_date, progress, labels, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		t.Title, t.Description, t.Status, t.Priority, t.StatusID, t.PriorityID,
		t.AssigneeID, t.CreatedBy, t.StartDate, t.DueDate, t.Progress,
		labelsJSON, t.CreatedAt, t.UpdatedAt).Scan(&t.ID)
}

func (r *taskRepo) GetByID(ctx context.Context, id uint) (*entity.Task, error) {
	t := &entity.Task{}
	var labelsJSON []byte
	query := `SELECT 
				t.id, t.title, t.description, t.status, t.priority, t.status_id, t.priority_id,
				t.assignee_id, t.created_by, t.start_date, t.due_date, t.progress, t.labels, t.created_at, t.updated_at,
				s.name as status_name, s.color as status_color, s.icon as status_icon,
				p.name as priority_name, p.color as priority_color, p.icon as priority_icon
			  FROM tasks t
			  LEFT JOIN categories s ON t.status_id = s.id
			  LEFT JOIN categories p ON t.priority_id = p.id
			  WHERE t.id = $1`

	var statusID, priorityID sql.NullInt64
	var statusName, statusColor, statusIcon sql.NullString
	var priorityName, priorityColor, priorityIcon sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &statusID, &priorityID,
		&t.AssigneeID, &t.CreatedBy, &t.StartDate, &t.DueDate, &t.Progress,
		&labelsJSON, &t.CreatedAt, &t.UpdatedAt,
		&statusName, &statusColor, &statusIcon,
		&priorityName, &priorityColor, &priorityIcon,
	)
	if err != nil {
		return nil, err
	}

	if statusID.Valid {
		t.StatusID = uint(statusID.Int64)
		t.StatusCat = &entity.Category{
			ID:    t.StatusID,
			Name:  statusName.String,
			Color: nullStringPtr(statusColor),
			Icon:  nullStringPtr(statusIcon),
		}
	}
	if priorityID.Valid {
		t.PriorityID = uint(priorityID.Int64)
		t.PriorityCat = &entity.Category{
			ID:    t.PriorityID,
			Name:  priorityName.String,
			Color: nullStringPtr(priorityColor),
			Icon:  nullStringPtr(priorityIcon),
		}
	}

	if err := json.Unmarshal(labelsJSON, &t.Labels); err != nil {
		t.Labels = []string{}
	}
	return t, nil
}

func (r *taskRepo) Update(ctx context.Context, t *entity.Task) error {
	labelsJSON, err := json.Marshal(t.Labels)
	if err != nil {
		labelsJSON = []byte("[]")
	}

	query := `UPDATE tasks SET
				title = $1,
				description = $2,
				status = $3,
				priority = $4,
				status_id = $5,
				priority_id = $6,
				assignee_id = $7,
				start_date = $8,
				due_date = $9,
				progress = $10,
				labels = $11,
				updated_at = $12
			  WHERE id = $13`

	_, err = r.db.ExecContext(ctx, query,
		t.Title, t.Description, t.Status, t.Priority, t.StatusID, t.PriorityID,
		t.AssigneeID, t.StartDate, t.DueDate, t.Progress,
		labelsJSON, time.Now(), t.ID)
	return err
}

func (r *taskRepo) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *taskRepo) ListByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Task, error) {
	query := `SELECT 
				t.id, t.title, t.description, t.status, t.priority, t.status_id, t.priority_id,
				t.assignee_id, t.created_by, t.start_date, t.due_date, t.progress, t.labels, t.created_at, t.updated_at,
				s.name as status_name, s.color as status_color, s.icon as status_icon,
				p.name as priority_name, p.color as priority_color, p.icon as priority_icon
			  FROM tasks t
			  LEFT JOIN categories s ON t.status_id = s.id
			  LEFT JOIN categories p ON t.priority_id = p.id
			  WHERE t.assignee_id = $1 OR t.created_by = $1 
			  ORDER BY t.created_at DESC
			  LIMIT $2 OFFSET $3`
	return r.list(ctx, query, userID, limit, offset)
}

func (r *taskRepo) CountByUser(ctx context.Context, userID uint) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM tasks WHERE assignee_id = $1 OR created_by = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *taskRepo) LogActivity(ctx context.Context, a *entity.TaskActivity) error {
	query := `INSERT INTO task_activities (task_id, actor_id, action, old_value, new_value, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	a.CreatedAt = time.Now()
	return r.db.QueryRowContext(ctx, query,
		a.TaskID, a.ActorID, a.Action, a.OldValue, a.NewValue, a.CreatedAt).Scan(&a.ID)
}

func (r *taskRepo) ListActivities(ctx context.Context, taskID uint) ([]entity.TaskActivity, error) {
	query := `SELECT id, task_id, actor_id, action, old_value, new_value, created_at
			  FROM task_activities WHERE task_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	activities := []entity.TaskActivity{}
	for rows.Next() {
		var a entity.TaskActivity
		if err := rows.Scan(&a.ID, &a.TaskID, &a.ActorID, &a.Action, &a.OldValue, &a.NewValue, &a.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, nil
}

func (r *taskRepo) list(ctx context.Context, query string, args ...interface{}) ([]entity.Task, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []entity.Task{}
	for rows.Next() {
		var t entity.Task
		var labelsJSON []byte
		var statusID, priorityID sql.NullInt64
		var statusName, statusColor, statusIcon sql.NullString
		var priorityName, priorityColor, priorityIcon sql.NullString

		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &statusID, &priorityID,
			&t.AssigneeID, &t.CreatedBy, &t.StartDate, &t.DueDate, &t.Progress,
			&labelsJSON, &t.CreatedAt, &t.UpdatedAt,
			&statusName, &statusColor, &statusIcon,
			&priorityName, &priorityColor, &priorityIcon,
		)
		if err != nil {
			return nil, err
		}

		if statusID.Valid {
			t.StatusID = uint(statusID.Int64)
			t.StatusCat = &entity.Category{
				ID:    t.StatusID,
				Name:  statusName.String,
				Color: nullStringPtr(statusColor),
				Icon:  nullStringPtr(statusIcon),
			}
		}
		if priorityID.Valid {
			t.PriorityID = uint(priorityID.Int64)
			t.PriorityCat = &entity.Category{
				ID:    t.PriorityID,
				Name:  priorityName.String,
				Color: nullStringPtr(priorityColor),
				Icon:  nullStringPtr(priorityIcon),
			}
		}

		if err := json.Unmarshal(labelsJSON, &t.Labels); err != nil {
			t.Labels = []string{}
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
