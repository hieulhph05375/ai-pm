package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ProjectRoleService interface {
	CreateRole(ctx context.Context, role *entity.ProjectRole) error
	GetRoleByID(ctx context.Context, id uint) (*entity.ProjectRole, error)
	GetRolesByProject(ctx context.Context, projectID uint) ([]entity.ProjectRole, error)
	UpdateRole(ctx context.Context, role *entity.ProjectRole) error
	DeleteRole(ctx context.Context, id uint) error

	GetAllPermissions(ctx context.Context) ([]entity.ProjectPermission, error)
	SetRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]entity.ProjectPermission, error)
}
