package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ProjectService interface {
	CreateProject(ctx context.Context, p *entity.Project) error
	GetProject(ctx context.Context, id int, userID int, isAdmin bool) (*entity.Project, error)
	UpdateProject(ctx context.Context, p *entity.Project) error
	DeleteProject(ctx context.Context, id int) error
	ListProjects(ctx context.Context, page, limit int, search string, status string, userID int, isAdmin bool) ([]entity.Project, int, error)
}
