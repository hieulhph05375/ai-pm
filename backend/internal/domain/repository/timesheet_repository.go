package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type TimesheetRepository interface {
	Create(ctx context.Context, timesheet *entity.Timesheet) error
	GetByID(ctx context.Context, id int) (*entity.Timesheet, error)
	Update(ctx context.Context, timesheet *entity.Timesheet) error
	Delete(ctx context.Context, id int) error
	ListByUser(ctx context.Context, userID int, limit, offset int) ([]entity.Timesheet, int, error)
	ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Timesheet, int, error)
}
