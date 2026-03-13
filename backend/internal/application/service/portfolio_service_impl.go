package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type portfolioService struct {
	repo repository.ProjectRepository
}

func NewPortfolioService(repo repository.ProjectRepository) PortfolioService {
	return &portfolioService{repo: repo}
}

func (s *portfolioService) GetPortfolioOverview(ctx context.Context, userID int, isAdmin bool) (*entity.PortfolioOverview, error) {
	return s.repo.GetPortfolioOverview(ctx, userID, isAdmin)
}
