package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

// StakeholderRepository defines the interface for stakeholder storage
type StakeholderRepository interface {
	// Global stakeholder operations
	Create(ctx context.Context, stakeholder *entity.Stakeholder) error
	GetByID(ctx context.Context, id int) (*entity.Stakeholder, error)
	Update(ctx context.Context, stakeholder *entity.Stakeholder) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, search string) ([]*entity.Stakeholder, error)
	ListWithPagination(ctx context.Context, search string, offset, limit int) ([]*entity.Stakeholder, int, error)
	Count(ctx context.Context, search string) (int, error)

	// Project stakeholder operations
	AssignToProject(ctx context.Context, projectID int, stakeholderID int, projectRole string, roleID uint) error
	UnassignFromProject(ctx context.Context, projectID int, stakeholderID int) error
	ListByProject(ctx context.Context, projectID int) ([]*entity.ProjectStakeholder, error)
}
