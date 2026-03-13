package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type stakeholderService struct {
	repo repository.StakeholderRepository
}

func NewStakeholderService(repo repository.StakeholderRepository) StakeholderService {
	return &stakeholderService{repo: repo}
}

func (s *stakeholderService) CreateStakeholder(ctx context.Context, stakeholder *entity.Stakeholder) error {
	return s.repo.Create(ctx, stakeholder)
}

func (s *stakeholderService) GetStakeholder(ctx context.Context, id int) (*entity.Stakeholder, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *stakeholderService) UpdateStakeholder(ctx context.Context, stakeholder *entity.Stakeholder) error {
	return s.repo.Update(ctx, stakeholder)
}

func (s *stakeholderService) DeleteStakeholder(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *stakeholderService) ListStakeholders(ctx context.Context, search string) ([]*entity.Stakeholder, error) {
	return s.repo.List(ctx, search)
}

func (s *stakeholderService) ListStakeholdersPaginated(ctx context.Context, search string, page, limit int) ([]*entity.Stakeholder, int, error) {
	offset := (page - 1) * limit
	return s.repo.ListWithPagination(ctx, search, offset, limit)
}

func (s *stakeholderService) AssignToProject(ctx context.Context, projectID int, stakeholderID int, role string, roleID uint) error {
	return s.repo.AssignToProject(ctx, projectID, stakeholderID, role, roleID)
}

func (s *stakeholderService) UnassignFromProject(ctx context.Context, projectID int, stakeholderID int) error {
	return s.repo.UnassignFromProject(ctx, projectID, stakeholderID)
}

func (s *stakeholderService) ListByProject(ctx context.Context, projectID int) ([]*entity.ProjectStakeholder, error) {
	return s.repo.ListByProject(ctx, projectID)
}
