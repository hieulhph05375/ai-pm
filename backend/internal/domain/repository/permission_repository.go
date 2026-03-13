package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type PermissionRepository interface {
	List(ctx context.Context) ([]entity.Permission, error)
	GetByID(ctx context.Context, id uint) (*entity.Permission, error)
	GetByName(ctx context.Context, name string) (*entity.Permission, error)
	Create(ctx context.Context, perm *entity.Permission) error
}
