package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type PortfolioService interface {
	GetPortfolioOverview(ctx context.Context, userID int, isAdmin bool) (*entity.PortfolioOverview, error)
}
