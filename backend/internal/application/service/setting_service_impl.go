package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{repo: repo}
}

func (s *settingService) Get(ctx context.Context, key string) (*entity.SystemSetting, error) {
	return s.repo.Get(ctx, key)
}

func (s *settingService) Set(ctx context.Context, key string, value any) error {
	setting := &entity.SystemSetting{
		Key:   key,
		Value: value,
	}
	return s.repo.Set(ctx, setting)
}

func (s *settingService) GetAll(ctx context.Context) ([]*entity.SystemSetting, error) {
	return s.repo.GetAll(ctx)
}
