package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ResourceRepository interface {
	GetWorkload(ctx context.Context, startDate, endDate string) (*entity.WorkloadOverview, error)
}
