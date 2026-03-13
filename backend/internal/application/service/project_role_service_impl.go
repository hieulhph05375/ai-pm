package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/infrastructure/db"
)

type projectRoleServiceImpl struct {
	repo db.ProjectRoleRepository
}

func NewProjectRoleService(repo db.ProjectRoleRepository) ProjectRoleService {
	return &projectRoleServiceImpl{repo: repo}
}

func (s *projectRoleServiceImpl) CreateRole(ctx context.Context, role *entity.ProjectRole) error {
	return s.repo.Create(ctx, role)
}

func (s *projectRoleServiceImpl) GetRoleByID(ctx context.Context, id uint) (*entity.ProjectRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectRoleServiceImpl) GetRolesByProject(ctx context.Context, projectID uint) ([]entity.ProjectRole, error) {
	return s.repo.GetByProject(ctx, projectID)
}

func (s *projectRoleServiceImpl) UpdateRole(ctx context.Context, role *entity.ProjectRole) error {
	return s.repo.Update(ctx, role)
}

func (s *projectRoleServiceImpl) DeleteRole(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *projectRoleServiceImpl) GetAllPermissions(ctx context.Context) ([]entity.ProjectPermission, error) {
	return s.repo.GetPermissions(ctx)
}

func (s *projectRoleServiceImpl) SetRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.repo.SetPermissions(ctx, roleID, permissionIDs)
}

func (s *projectRoleServiceImpl) GetRolePermissions(ctx context.Context, roleID uint) ([]entity.ProjectPermission, error) {
	return s.repo.GetRolePermissions(ctx, roleID)
}
