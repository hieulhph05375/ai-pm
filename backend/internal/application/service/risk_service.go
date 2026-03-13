package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type RiskService interface {
	CreateRisk(ctx context.Context, r *entity.Risk) error
	ListRisks(ctx context.Context, projectID int, limit, offset int) ([]entity.Risk, int, error)
	UpdateRisk(ctx context.Context, r *entity.Risk) error
	DeleteRisk(ctx context.Context, id int) error
}
