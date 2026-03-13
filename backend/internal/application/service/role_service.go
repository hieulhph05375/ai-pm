package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type RoleService interface {
	ListRoles(ctx context.Context) ([]entity.Role, error)
	GetRole(ctx context.Context, id uint) (*entity.Role, error)
	GetRoleWithPermissions(ctx context.Context, id uint) (*entity.RoleWithPermissions, error)
	CreateRole(ctx context.Context, role *entity.Role) error
	UpdateRole(ctx context.Context, role *entity.Role) error
	DeleteRole(ctx context.Context, id uint) error

	AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
	ListPermissions(ctx context.Context) ([]entity.Permission, error)
}

type roleService struct {
	roleRepo repository.RoleRepository
	permRepo repository.PermissionRepository
}

func NewRoleService(rr repository.RoleRepository, pr repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo: rr,
		permRepo: pr,
	}
}

func (s *roleService) ListRoles(ctx context.Context) ([]entity.Role, error) {
	return s.roleRepo.List(ctx)
}

func (s *roleService) GetRole(ctx context.Context, id uint) (*entity.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

func (s *roleService) GetRoleWithPermissions(ctx context.Context, id uint) (*entity.RoleWithPermissions, error) {
	role, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	perms, err := s.roleRepo.GetPermissionsByRoleID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.RoleWithPermissions{
		Role:        *role,
		Permissions: perms,
	}, nil
}

func (s *roleService) CreateRole(ctx context.Context, role *entity.Role) error {
	return s.roleRepo.Create(ctx, role)
}

func (s *roleService) UpdateRole(ctx context.Context, role *entity.Role) error {
	return s.roleRepo.Update(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepo.Delete(ctx, id)
}

func (s *roleService) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(ctx, roleID, permissionIDs)
}

func (s *roleService) ListPermissions(ctx context.Context) ([]entity.Permission, error) {
	return s.permRepo.List(ctx)
}
