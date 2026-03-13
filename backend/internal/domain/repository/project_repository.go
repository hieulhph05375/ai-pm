package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ProjectRepository interface {
	Create(ctx context.Context, p *entity.Project) error
	GetByID(ctx context.Context, id int, userID int, isAdmin bool) (*entity.Project, error)
	GetByProjectID(ctx context.Context, projectID string, userID int, isAdmin bool) (*entity.Project, error)
	Update(ctx context.Context, p *entity.Project) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int, search string, status string, userID int, isAdmin bool) ([]entity.Project, int, error)
	GetPortfolioOverview(ctx context.Context, userID int, isAdmin bool) (*entity.PortfolioOverview, error)
}
