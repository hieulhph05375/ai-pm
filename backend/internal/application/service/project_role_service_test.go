package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"testing"
)

type mockProjectRoleRepo struct {
	roles       map[uint]*entity.ProjectRole
	permissions map[uint][]entity.ProjectPermission
}

func (m *mockProjectRoleRepo) Create(ctx context.Context, role *entity.ProjectRole) error {
	role.ID = uint(len(m.roles) + 1)
	m.roles[role.ID] = role
	return nil
}

func (m *mockProjectRoleRepo) GetByID(ctx context.Context, id uint) (*entity.ProjectRole, error) {
	if role, ok := m.roles[id]; ok {
		return role, nil
	}
	return nil, nil // Normally return an error, returning nil for simplicity
}

func (m *mockProjectRoleRepo) GetByProject(ctx context.Context, projectID uint) ([]entity.ProjectRole, error) {
	var res []entity.ProjectRole
	for _, role := range m.roles {
		if role.ProjectID == projectID {
			res = append(res, *role)
		}
	}
	return res, nil
}

func (m *mockProjectRoleRepo) Update(ctx context.Context, role *entity.ProjectRole) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockProjectRoleRepo) Delete(ctx context.Context, id uint) error {
	delete(m.roles, id)
	return nil
}

func (m *mockProjectRoleRepo) GetPermissions(ctx context.Context) ([]entity.ProjectPermission, error) {
	return []entity.ProjectPermission{
		{ID: 1, Name: "project:wbs:create", Description: "Create WBS"},
		{ID: 2, Name: "project:wbs:update", Description: "Update WBS"},
		{ID: 3, Name: "project:wbs:delete", Description: "Delete WBS"},
	}, nil
}

func (m *mockProjectRoleRepo) SetPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	var perms []entity.ProjectPermission
	for _, pid := range permissionIDs {
		perms = append(perms, entity.ProjectPermission{ID: pid})
	}
	m.permissions[roleID] = perms
	return nil
}

func (m *mockProjectRoleRepo) GetRolePermissions(ctx context.Context, roleID uint) ([]entity.ProjectPermission, error) {
	return m.permissions[roleID], nil
}

func TestProjectRoleService(t *testing.T) {
	repo := &mockProjectRoleRepo{
		roles:       make(map[uint]*entity.ProjectRole),
		permissions: make(map[uint][]entity.ProjectPermission),
	}
	service := NewProjectRoleService(repo)
	ctx := context.Background()

	t.Run("Create Role", func(t *testing.T) {
		role := &entity.ProjectRole{
			ProjectID:   1,
			Name:        "Test Role",
			Description: "Test Description",
		}
		err := service.CreateRole(ctx, role)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if role.ID == 0 {
			t.Errorf("Expected assigned ID, got 0")
		}
	})

	t.Run("Get Roles By Project", func(t *testing.T) {
		roles, err := service.GetRolesByProject(ctx, 1)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if len(roles) != 1 {
			t.Errorf("Expected 1 role, got %d", len(roles))
		}
	})

	t.Run("Set and Get Role Permissions", func(t *testing.T) {
		// role.ID from first test should be 1
		err := service.SetRolePermissions(ctx, 1, []uint{1, 2})
		if err != nil {
			t.Errorf("Expected nil error setting permissions, got %v", err)
		}

		perms, err := service.GetRolePermissions(ctx, 1)
		if err != nil {
			t.Errorf("Expected nil error getting permissions, got %v", err)
		}
		if len(perms) != 2 {
			t.Errorf("Expected 2 permissions, got %d", len(perms))
		}
	})

	t.Run("Delete Role", func(t *testing.T) {
		err := service.DeleteRole(ctx, 1)
		if err != nil {
			t.Errorf("Expected nil error deleting role, got %v", err)
		}

		roles, err := service.GetRolesByProject(ctx, 1)
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}
		if len(roles) != 0 {
			t.Errorf("Expected empty roles map, got %d", len(roles))
		}
	})
}
