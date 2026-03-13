package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uint) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uint) error
	// ListByUser returns tasks assigned to a specific user with pagination
	ListByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Task, error)
	// CountByUser returns total number of tasks for a user
	CountByUser(ctx context.Context, userID uint) (int, error)
	// LogActivity records an activity entry for audit trail
	LogActivity(ctx context.Context, activity *entity.TaskActivity) error
	// ListActivities returns activity log for a task
	ListActivities(ctx context.Context, taskID uint) ([]entity.TaskActivity, error)
}
