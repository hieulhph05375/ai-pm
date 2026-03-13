package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type ProjectMemberService interface {
	AddMember(ctx context.Context, member *entity.ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID int) error
	UpdateMemberRole(ctx context.Context, projectID, userID int, roleID int) error
	GetMembersByProject(ctx context.Context, projectID int, page, limit int) ([]entity.ProjectMember, int, error)
	GetProjectsByUser(ctx context.Context, userID int) ([]entity.ProjectMember, error)
	IsMember(ctx context.Context, projectID, userID int) (bool, int, error)
}

type projectMemberServiceImpl struct {
	repo repository.ProjectMemberRepository
}

func NewProjectMemberService(repo repository.ProjectMemberRepository) ProjectMemberService {
	return &projectMemberServiceImpl{repo: repo}
}

func (s *projectMemberServiceImpl) AddMember(ctx context.Context, member *entity.ProjectMember) error {
	return s.repo.AddMember(ctx, member)
}

func (s *projectMemberServiceImpl) RemoveMember(ctx context.Context, projectID, userID int) error {
	return s.repo.RemoveMember(ctx, projectID, userID)
}

func (s *projectMemberServiceImpl) UpdateMemberRole(ctx context.Context, projectID, userID int, roleID int) error {
	return s.repo.UpdateMemberRole(ctx, projectID, userID, roleID)
}

func (s *projectMemberServiceImpl) GetMembersByProject(ctx context.Context, projectID int, page, limit int) ([]entity.ProjectMember, int, error) {
	return s.repo.GetMembersByProject(ctx, projectID, page, limit)
}

func (s *projectMemberServiceImpl) GetProjectsByUser(ctx context.Context, userID int) ([]entity.ProjectMember, error) {
	return s.repo.GetProjectsByUser(ctx, userID)
}

func (s *projectMemberServiceImpl) IsMember(ctx context.Context, projectID, userID int) (bool, int, error) {
	return s.repo.IsMember(ctx, projectID, userID)
}
