package service

import (
	"context"
	"project-mgmt/backend/internal/domain/entity"
)

type SnapshotService interface {
	CaptureAllProjectsSnapshot(ctx context.Context) error
	GetProjectTrends(ctx context.Context, projectID int) ([]entity.ProjectSnapshot, error)
	GetMilestoneTrends(ctx context.Context, projectID int) ([]entity.MilestoneSnapshot, error)
}
