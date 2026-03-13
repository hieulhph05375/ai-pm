package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type ProjectMemberRepository interface {
	AddMember(ctx context.Context, member *entity.ProjectMember) error
	RemoveMember(ctx context.Context, projectID, userID int) error
	UpdateMemberRole(ctx context.Context, projectID, userID int, roleID int) error
	GetMembersByProject(ctx context.Context, projectID int, page, limit int) ([]entity.ProjectMember, int, error)
	GetProjectsByUser(ctx context.Context, userID int) ([]entity.ProjectMember, error)
	IsMember(ctx context.Context, projectID, userID int) (bool, int, error)
}
