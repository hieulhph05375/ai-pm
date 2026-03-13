package service

import (
	"context"
	"fmt"
	"math"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strings"
)

type snapshotService struct {
	projectRepo  repository.ProjectRepository
	wbsRepo      repository.WBSRepository
	snapshotRepo repository.SnapshotRepository
}

func NewSnapshotService(
	pRepo repository.ProjectRepository,
	wRepo repository.WBSRepository,
	sRepo repository.SnapshotRepository,
) SnapshotService {
	return &snapshotService{
		projectRepo:  pRepo,
		wbsRepo:      wRepo,
		snapshotRepo: sRepo,
	}
}

func (s *snapshotService) CaptureAllProjectsSnapshot(ctx context.Context) error {
	// List all projects (simple list for snapshotting)
	// Using offset=0, limit=1000 to get a large set of projects
	projects, _, err := s.projectRepo.List(ctx, 0, 1000, "", "", 0, true)
	if err != nil {
		return err
	}

	for _, project := range projects {
		// 1. Calculate PMI Metrics for the project
		nodes, _, err := s.wbsRepo.GetProjectTree(ctx, project.ID, entity.WBSFilter{})
		if err != nil {
			fmt.Printf("DEBUG: GetProjectTree failed for project %d: %v\n", project.ID, err)
			continue // Skip failed projects
		}

		var totalPV, totalEV, totalAC float64
		isParent := make(map[string]bool)
		for _, node := range nodes {
			for i := len(nodes) - 1; i >= 0; i-- {
				if nodes[i].Path != node.Path && strings.HasPrefix(nodes[i].Path, node.Path+".") {
					isParent[node.Path] = true
					break
				}
			}
		}

		for _, node := range nodes {
			if !isParent[node.Path] {
				totalPV += node.PlannedValue
				totalEV += (float64(node.Progress) / 100.0) * node.PlannedValue
				totalAC += node.ActualCost
			}
		}

		spi := 1.0
		if totalPV > 0 {
			spi = totalEV / totalPV
		}
		cpi := 1.0
		if totalAC > 0 {
			cpi = totalEV / totalAC
		}

		// Precision handling
		spi = math.Round(spi*100) / 100
		cpi = math.Round(cpi*100) / 100

		// 2. Save Project Snapshot
		ps := &entity.ProjectSnapshot{
			ProjectID: project.ID,
			SPI:       spi,
			CPI:       cpi,
			EV:        totalEV,
			AC:        totalAC,
			PV:        totalPV,
			Progress:  project.Progress,
		}
		if err := s.snapshotRepo.CreateProjectSnapshot(ctx, ps); err != nil {
			fmt.Printf("DEBUG: CreateProjectSnapshot failed for project %d: %v\n", project.ID, err)
		}

		// 3. Save Milestone Snapshots
		for _, node := range nodes {
			if node.Type == entity.TypeMilestone {
				if node.PlannedEndDate == nil {
					continue
				}
				ms := &entity.MilestoneSnapshot{
					ProjectID:     project.ID,
					NodeID:        node.ID,
					MilestoneName: node.Title,
					PlannedDate:   *node.PlannedEndDate,
				}
				// If progress is 100%, consider PlannedEndDate as actual date
				if node.Progress == 100 {
					actual := *node.PlannedEndDate
					ms.ActualDate = &actual
				}
				if err := s.snapshotRepo.CreateMilestoneSnapshot(ctx, ms); err != nil {
					fmt.Printf("DEBUG: CreateMilestoneSnapshot failed for node %d: %v\n", node.ID, err)
				}
			}
		}
	}

	return nil
}

func (s *snapshotService) GetProjectTrends(ctx context.Context, projectID int) ([]entity.ProjectSnapshot, error) {
	return s.snapshotRepo.GetProjectSnapshots(ctx, projectID)
}

func (s *snapshotService) GetMilestoneTrends(ctx context.Context, projectID int) ([]entity.MilestoneSnapshot, error) {
	return s.snapshotRepo.GetMilestoneSnapshots(ctx, projectID)
}
