package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type RoleRepository interface {
	List(ctx context.Context) ([]entity.Role, error)
	GetByID(ctx context.Context, id uint) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	Create(ctx context.Context, role *entity.Role) error
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id uint) error

	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]entity.Permission, error)
}
