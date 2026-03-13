package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type IssueRepository interface {
	Create(ctx context.Context, i *entity.Issue) error
	ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Issue, error)
	CountByProject(ctx context.Context, projectID int) (int, error)
	Update(ctx context.Context, i *entity.Issue) error
	Delete(ctx context.Context, id int) error
}
