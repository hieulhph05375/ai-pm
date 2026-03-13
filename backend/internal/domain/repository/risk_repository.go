package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type RiskRepository interface {
	Create(ctx context.Context, r *entity.Risk) error
	ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Risk, error)
	CountByProject(ctx context.Context, projectID int) (int, error)
	Update(ctx context.Context, r *entity.Risk) error
	Delete(ctx context.Context, id int) error
}
