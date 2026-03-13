package service

import (
	"context"
	"errors"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *projectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(ctx context.Context, p *entity.Project) error {
	existing, err := s.repo.GetByProjectID(ctx, p.ProjectID, 0, true)
	if err == nil && existing != nil {
		return errors.New("mã dự án đã tồn tại trên hệ thống")
	}
	return s.repo.Create(ctx, p)
}

func (s *projectService) GetByProjectID(ctx context.Context, projectID string, userID int, isAdmin bool) (*entity.Project, error) {
	return s.repo.GetByProjectID(ctx, projectID, userID, isAdmin)
}

func (s *projectService) GetProject(ctx context.Context, id int, userID int, isAdmin bool) (*entity.Project, error) {
	return s.repo.GetByID(ctx, id, userID, isAdmin)
}

func (s *projectService) UpdateProject(ctx context.Context, p *entity.Project) error {
	existing, err := s.repo.GetByProjectID(ctx, p.ProjectID, 0, true)
	if err == nil && existing != nil && existing.ID != p.ID {
		return errors.New("mã dự án đã tồn tại trên hệ thống")
	}
	return s.repo.Update(ctx, p)
}

func (s *projectService) DeleteProject(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *projectService) ListProjects(ctx context.Context, page, limit int, search string, status string, userID int, isAdmin bool) ([]entity.Project, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.List(ctx, offset, limit, search, status, userID, isAdmin)
}
