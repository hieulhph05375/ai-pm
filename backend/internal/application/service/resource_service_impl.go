package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type resourceService struct {
	repo repository.ResourceRepository
}

func NewResourceService(repo repository.ResourceRepository) ResourceService {
	return &resourceService{repo: repo}
}

func (s *resourceService) GetWorkload(ctx context.Context, startDate, endDate string) (*entity.WorkloadOverview, error) {
	return s.repo.GetWorkload(ctx, startDate, endDate)
}
