package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"time"
)

// HolidayRepository defines the interface for holiday storage
type HolidayRepository interface {
	Create(ctx context.Context, holiday *entity.Holiday) error
	GetByID(ctx context.Context, id int) (*entity.Holiday, error)
	GetByDate(ctx context.Context, date time.Time) (*entity.Holiday, error)
	Update(ctx context.Context, holiday *entity.Holiday) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error)
	ListWithPagination(ctx context.Context, start, end time.Time, offset, limit int) ([]*entity.Holiday, int, error)
	Count(ctx context.Context, start, end time.Time) (int, error)
}
