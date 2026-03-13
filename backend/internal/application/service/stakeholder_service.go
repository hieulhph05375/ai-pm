package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type StakeholderService interface {
	CreateStakeholder(ctx context.Context, s *entity.Stakeholder) error
	GetStakeholder(ctx context.Context, id int) (*entity.Stakeholder, error)
	UpdateStakeholder(ctx context.Context, s *entity.Stakeholder) error
	DeleteStakeholder(ctx context.Context, id int) error
	ListStakeholders(ctx context.Context, search string) ([]*entity.Stakeholder, error)
	ListStakeholdersPaginated(ctx context.Context, search string, page, limit int) ([]*entity.Stakeholder, int, error)

	AssignToProject(ctx context.Context, projectID int, stakeholderID int, role string, roleID uint) error
	UnassignFromProject(ctx context.Context, projectID int, stakeholderID int) error
	ListByProject(ctx context.Context, projectID int) ([]*entity.ProjectStakeholder, error)
}
