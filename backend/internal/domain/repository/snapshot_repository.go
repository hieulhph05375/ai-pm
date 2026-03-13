package repository

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type SnapshotRepository interface {
	CreateProjectSnapshot(ctx context.Context, snapshot *entity.ProjectSnapshot) error
	CreateMilestoneSnapshot(ctx context.Context, snapshot *entity.MilestoneSnapshot) error
	GetProjectSnapshots(ctx context.Context, projectID int) ([]entity.ProjectSnapshot, error)
	GetMilestoneSnapshots(ctx context.Context, projectID int) ([]entity.MilestoneSnapshot, error)
}
