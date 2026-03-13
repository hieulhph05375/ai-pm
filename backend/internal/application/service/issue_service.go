package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type IssueService interface {
	CreateIssue(ctx context.Context, i *entity.Issue) error
	ListIssues(ctx context.Context, projectID int, limit, offset int) ([]entity.Issue, int, error)
	UpdateIssue(ctx context.Context, i *entity.Issue) error
	DeleteIssue(ctx context.Context, id int) error
}
