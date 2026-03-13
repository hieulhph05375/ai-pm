package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"time"
)

type holidayService struct {
	repo repository.HolidayRepository
}

func NewHolidayService(repo repository.HolidayRepository) HolidayService {
	return &holidayService{repo: repo}
}

func (s *holidayService) CreateHoliday(ctx context.Context, holiday *entity.Holiday) error {
	return s.repo.Create(ctx, holiday)
}

func (s *holidayService) GetHoliday(ctx context.Context, id int) (*entity.Holiday, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *holidayService) UpdateHoliday(ctx context.Context, holiday *entity.Holiday) error {
	return s.repo.Update(ctx, holiday)
}

func (s *holidayService) DeleteHoliday(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *holidayService) ListHolidays(ctx context.Context, start, end time.Time) ([]*entity.Holiday, error) {
	return s.repo.List(ctx, start, end)
}

func (s *holidayService) ListHolidaysPaginated(ctx context.Context, start, end time.Time, page, limit int) ([]*entity.Holiday, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.ListWithPagination(ctx, start, end, offset, limit)
}
