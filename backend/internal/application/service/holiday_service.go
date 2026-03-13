package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"time"
)

type HolidayService interface {
	CreateHoliday(ctx context.Context, h *entity.Holiday) error
	GetHoliday(ctx context.Context, id int) (*entity.Holiday, error)
	UpdateHoliday(ctx context.Context, h *entity.Holiday) error
	DeleteHoliday(ctx context.Context, id int) error
	ListHolidays(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error)
	ListHolidaysPaginated(ctx context.Context, start, end time.Time, page, limit int) ([]*entity.Holiday, int, error)
}
