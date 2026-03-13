package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type issueService struct {
	repo repository.IssueRepository
}

func NewIssueService(repo repository.IssueRepository) IssueService {
	return &issueService{repo: repo}
}

func (s *issueService) CreateIssue(ctx context.Context, i *entity.Issue) error {
	// Validation can be added here once we have a clear source of default IDs
	return s.repo.Create(ctx, i)
}

func (s *issueService) ListIssues(ctx context.Context, projectID int, limit, offset int) ([]entity.Issue, int, error) {
	issues, err := s.repo.ListByProject(ctx, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByProject(ctx, projectID)
	if err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (s *issueService) UpdateIssue(ctx context.Context, i *entity.Issue) error {
	return s.repo.Update(ctx, i)
}

func (s *issueService) DeleteIssue(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
