package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ResourceService interface {
	GetWorkload(ctx context.Context, startDate, endDate string) (*entity.WorkloadOverview, error)
}
