package service

import (
	"context"
	"math"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strings"
	"time"
)

type PMIStats struct {
	PV          float64 `json:"pv"`  // Planned Value
	EV          float64 `json:"ev"`  // Earned Value
	AC          float64 `json:"ac"`  // Actual Cost
	SPI         float64 `json:"spi"` // Schedule Performance Index
	CPI         float64 `json:"cpi"` // Cost Performance Index
	EAC         float64 `json:"eac"` // Estimate At Completion
	NeedsUpdate bool    `json:"needs_update"`
}

type ReportService interface {
	GetProjectPMIStats(ctx context.Context, projectID int, userID int, isAdmin bool) (*PMIStats, error)
}

type reportService struct {
	projectRepo repository.ProjectRepository
	wbsRepo     repository.WBSRepository
}

func NewReportService(pRepo repository.ProjectRepository, wRepo repository.WBSRepository) ReportService {
	return &reportService{
		projectRepo: pRepo,
		wbsRepo:     wRepo,
	}
}

func (s *reportService) GetProjectPMIStats(ctx context.Context, projectID int, userID int, isAdmin bool) (*PMIStats, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID, userID, isAdmin)
	if err != nil {
		return nil, err
	}

	filter := entity.WBSFilter{
		FetchAll: true,
		Fields:   []string{"path", "planned_value", "progress", "actual_cost"},
	} // fetch the entire WBS tree unpaginated, only specific fields
	nodes, _, err := s.wbsRepo.GetProjectTree(ctx, projectID, filter)
	if err != nil {
		return nil, err
	}

	stats := &PMIStats{}
	var totalPV, totalEV, totalAC float64

	// Map to track which nodes are parents (efficient O(N) approach)
	isParent := make(map[string]bool)
	for _, node := range nodes {
		lastDot := strings.LastIndex(node.Path, ".")
		if lastDot != -1 {
			parentPath := node.Path[:lastDot]
			isParent[parentPath] = true
		}
	}

	for _, node := range nodes {
		// Only sum leaf nodes (nodes that are not parents)
		if !isParent[node.Path] {
			totalPV += node.PlannedValue
			totalEV += (node.Progress / 100.0) * node.PlannedValue
			totalAC += node.ActualCost
		}
	}

	stats.PV = totalPV
	stats.EV = totalEV
	stats.AC = totalAC

	if totalPV > 0 {
		stats.SPI = totalEV / totalPV
	} else {
		stats.SPI = 1.0
	}

	if totalAC > 0 {
		stats.CPI = totalEV / totalAC
	} else {
		stats.CPI = 1.0
	}

	// EAC = BAC / CPI (Simplified)
	if stats.CPI > 0 {
		stats.EAC = totalPV / stats.CPI
	} else {
		stats.EAC = totalPV
	}

	// Logic for Reminder
	// Seven days threshold
	threshold := 7 * 24 * time.Hour
	if project.LastReminderAt != nil {
		if time.Since(*project.LastReminderAt) > threshold {
			stats.NeedsUpdate = true
		}
	} else if time.Since(project.UpdatedAt) > threshold {
		stats.NeedsUpdate = true
	}

	// Precision handling
	stats.SPI = math.Round(stats.SPI*100) / 100
	stats.CPI = math.Round(stats.CPI*100) / 100
	stats.EAC = math.Round(stats.EAC*100) / 100

	return stats, nil
}
