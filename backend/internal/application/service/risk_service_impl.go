package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type riskService struct {
	repo repository.RiskRepository
}

func NewRiskService(repo repository.RiskRepository) RiskService {
	return &riskService{repo: repo}
}

func (s *riskService) CreateRisk(ctx context.Context, r *entity.Risk) error {
	if r.Status == "" {
		r.Status = "Open"
	}
	return s.repo.Create(ctx, r)
}

func (s *riskService) ListRisks(ctx context.Context, projectID int, limit, offset int) ([]entity.Risk, int, error) {
	risks, err := s.repo.ListByProject(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByProject(ctx, projectID)
	if err != nil {
		return nil, 0, err
	}
	return risks, total, nil
}

func (s *riskService) UpdateRisk(ctx context.Context, r *entity.Risk) error {
	return s.repo.Update(ctx, r)
}

func (s *riskService) DeleteRisk(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
