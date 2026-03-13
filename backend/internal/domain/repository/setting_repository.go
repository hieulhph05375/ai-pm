package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type SettingRepository interface {
	Get(ctx context.Context, key string) (*entity.SystemSetting, error)
	Set(ctx context.Context, setting *entity.SystemSetting) error
	GetAll(ctx context.Context) ([]*entity.SystemSetting, error)
}
