package service

import (
	"context"
	"errors"
	"fmt"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
	"strings"
)

func validateNodeDates(node *entity.WBSNode, parent *entity.WBSNode) error {
	if node.PlannedStartDate != nil && parent.PlannedStartDate != nil {
		if node.PlannedStartDate.Before(*parent.PlannedStartDate) {
			return fmt.Errorf("node planned start date (%s) cannot be earlier than parent's planned start date (%s)",
				node.PlannedStartDate.Format("2006-01-02"), parent.PlannedStartDate.Format("2006-01-02"))
		}
	}
	if node.PlannedEndDate != nil && parent.PlannedEndDate != nil {
		if node.PlannedEndDate.After(*parent.PlannedEndDate) {
			return fmt.Errorf("node planned end date (%s) cannot be later than parent's planned end date (%s)",
				node.PlannedEndDate.Format("2006-01-02"), parent.PlannedEndDate.Format("2006-01-02"))
		}
	}
	if node.PlannedStartDate != nil && node.PlannedEndDate != nil {
		if node.PlannedStartDate.After(*node.PlannedEndDate) {
			return errors.New("node planned start date cannot be after planned end date")
		}
	}
	return nil
}

type WBSService interface {
	GetProjectTree(ctx context.Context, projectID int, filter entity.WBSFilter) ([]entity.WBSNode, int, error)
	CreateNode(ctx context.Context, node *entity.WBSNode, parentPath string) error
	UpdateNode(ctx context.Context, node *entity.WBSNode) error
	DeleteNode(ctx context.Context, id int) error
	GetNodeByID(ctx context.Context, id int) (*entity.WBSNode, error)
	ListDependencies(ctx context.Context, projectID int) ([]entity.WBSDependency, error)
	CreateDependency(ctx context.Context, dep *entity.WBSDependency) error
	DeleteDependency(ctx context.Context, depID int) error

	// Comment methods
	AddComment(ctx context.Context, comment *entity.WBSComment) error
	ListComments(ctx context.Context, nodeID int) ([]entity.WBSComment, error)
	DeleteComment(ctx context.Context, commentID int) error
	UpdateComment(ctx context.Context, commentID int, content string) error
	GetCommentByID(ctx context.Context, commentID int) (*entity.WBSComment, error)

	// Baseline methods
	CreateBaseline(ctx context.Context, projectID int, name string, description string, userID int) (*entity.WBSBaseline, error)
	GetBaselines(ctx context.Context, projectID int) ([]entity.WBSBaseline, error)
	GetBaselineNodes(ctx context.Context, baselineID int) ([]entity.WBSBaselineNode, error)
}

type wbsService struct {
	repo repository.WBSRepository
}

func NewWBSService(repo repository.WBSRepository) WBSService {
	return &wbsService{repo: repo}
}

func (s *wbsService) GetProjectTree(ctx context.Context, projectID int, filter entity.WBSFilter) ([]entity.WBSNode, int, error) {
	return s.repo.GetProjectTree(ctx, projectID, filter)
}

func (s *wbsService) CreateNode(ctx context.Context, node *entity.WBSNode, parentPath string) error {
	// Validate child dates
	if node.PlannedStartDate != nil && node.PlannedEndDate != nil {
		if node.PlannedStartDate.After(*node.PlannedEndDate) {
			return errors.New("node planned start date cannot be after planned end date")
		}
	}

	// Validate Depth limits (5 levels max)
	if parentPath != "" {
		parts := strings.Split(parentPath, ".")
		if len(parts) >= 5 {
			return errors.New("cannot create node: exceeds maximum WBS depth of 5")
		}

		// Retrieve parent for date validation
		parent, err := s.repo.GetNodeByPath(ctx, node.ProjectID, parentPath)
		if err == nil && parent != nil {
			if err := validateNodeDates(node, parent); err != nil {
				return err
			}
		}
	}

	err := s.repo.CreateNode(ctx, node, parentPath)
	if err != nil {
		return err
	}

	// Trigger Roll Up
	if parentPath != "" {
		return s.applyRollup(ctx, node.ProjectID, parentPath)
	}

	return nil
}

func (s *wbsService) UpdateNode(ctx context.Context, node *entity.WBSNode) error {
	// 1. Fetch existing node to get Path and ProjectID
	existing, err := s.repo.GetNodeByID(ctx, node.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch node for update: %w", err)
	}

	// 2. Prepare node for update (keep metadata from existing)
	node.Path = existing.Path
	node.ProjectID = existing.ProjectID

	// 3. Validate base dates
	if node.PlannedStartDate != nil && node.PlannedEndDate != nil {
		if node.PlannedStartDate.After(*node.PlannedEndDate) {
			return errors.New("node planned start date cannot be after planned end date")
		}
	}

	// 4. Validate against parent
	if node.Path != "" && strings.Contains(node.Path, ".") {
		parts := strings.Split(node.Path, ".")
		parentPath := strings.Join(parts[:len(parts)-1], ".")
		parent, err := s.repo.GetNodeByPath(ctx, node.ProjectID, parentPath)
		if err == nil && parent != nil {
			if err := validateNodeDates(node, parent); err != nil {
				return err
			}
		}
	}

	// 5. Validate against children
	if node.Path != "" {
		children, err := s.repo.GetChildren(ctx, node.Path)
		if err == nil {
			for i := range children {
				if err := validateNodeDates(&children[i], node); err != nil {
					return fmt.Errorf("update rejected: child node '%s' would be out of bounds. %w", children[i].Title, err)
				}
			}
		}
	}

	// 6. Persist update
	err = s.repo.UpdateNode(ctx, node)
	if err != nil {
		return err
	}

	// 7. Trigger Roll Up
	if node.Path != "" && strings.Contains(node.Path, ".") {
		parts := strings.Split(node.Path, ".")
		parentPath := strings.Join(parts[:len(parts)-1], ".")
		return s.applyRollup(ctx, node.ProjectID, parentPath)
	}

	return nil
}

func (s *wbsService) DeleteNode(ctx context.Context, id int) error {
	node, err := s.repo.GetNodeByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteNode(ctx, id)
	if err != nil {
		return err
	}

	// Trigger Roll Up
	if node.Path != "" && strings.Contains(node.Path, ".") {
		parts := strings.Split(node.Path, ".")
		parentPath := strings.Join(parts[:len(parts)-1], ".")
		return s.applyRollup(ctx, node.ProjectID, parentPath)
	}

	return nil
}

func (s *wbsService) GetNodeByID(ctx context.Context, id int) (*entity.WBSNode, error) {
	return s.repo.GetNodeByID(ctx, id)
}

func (s *wbsService) AddComment(ctx context.Context, comment *entity.WBSComment) error {
	return s.repo.AddComment(ctx, comment)
}

func (s *wbsService) ListComments(ctx context.Context, nodeID int) ([]entity.WBSComment, error) {
	return s.repo.ListComments(ctx, nodeID)
}

func (s *wbsService) DeleteComment(ctx context.Context, commentID int) error {
	return s.repo.DeleteComment(ctx, commentID)
}

func (s *wbsService) UpdateComment(ctx context.Context, commentID int, content string) error {
	return s.repo.UpdateComment(ctx, commentID, content)
}

func (s *wbsService) GetCommentByID(ctx context.Context, commentID int) (*entity.WBSComment, error) {
	return s.repo.GetCommentByID(ctx, commentID)
}

func (s *wbsService) ListDependencies(ctx context.Context, projectID int) ([]entity.WBSDependency, error) {
	return s.repo.ListDependencies(ctx, projectID)
}

func (s *wbsService) CreateDependency(ctx context.Context, dep *entity.WBSDependency) error {
	return s.repo.CreateDependency(ctx, dep)
}

func (s *wbsService) DeleteDependency(ctx context.Context, depID int) error {
	return s.repo.DeleteDependency(ctx, depID)
}

func (s *wbsService) CreateBaseline(ctx context.Context, projectID int, name string, description string, userID int) (*entity.WBSBaseline, error) {
	// Let's verify the project has at least some nodes before baselining it.
	filter := entity.WBSFilter{} // No filters
	nodes, _, err := s.repo.GetProjectTree(ctx, projectID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch project tree for baseline: %w", err)
	}

	if len(nodes) == 0 {
		return nil, errors.New("cannot create baseline for an empty project")
	}

	// Create baseline entity
	baseline := &entity.WBSBaseline{
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		CreatedBy:   userID,
	}

	// Save baseline metadata
	err = s.repo.CreateBaseline(ctx, baseline)
	if err != nil {
		return nil, fmt.Errorf("failed to save baseline metadata to database: %w", err)
	}

	// Trigger bulk SQL copy
	err = s.repo.CopyNodesToBaseline(ctx, projectID, baseline.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to copy wbs nodes to baseline: %w", err)
	}

	return baseline, nil
}

func (s *wbsService) GetBaselines(ctx context.Context, projectID int) ([]entity.WBSBaseline, error) {
	return s.repo.GetBaselines(ctx, projectID)
}

func (s *wbsService) GetBaselineNodes(ctx context.Context, baselineID int) ([]entity.WBSBaselineNode, error) {
	return s.repo.GetBaselineNodes(ctx, baselineID)
}

// applyRollup recursively updates parent node metrics (progress, PV, AC, effort)
func (s *wbsService) applyRollup(ctx context.Context, projectID int, parentPath string) error {
	if parentPath == "" {
		return nil
	}

	// 1. Get the parent node
	parent, err := s.repo.GetNodeByPath(ctx, projectID, parentPath)
	if err != nil {
		return err
	}

	// 2. Get all immediate children
	children, err := s.repo.GetImmediateChildren(ctx, projectID, parentPath)
	if err != nil {
		return err
	}

	if len(children) == 0 {
		return nil // Should not happen if we're rolling up from a child update
	}

	// 3. Calculate roll-up metrics
	var totalPV, totalEV, totalAC, totalEstimatedEffort, totalActualEffort float64
	for _, child := range children {
		totalPV += child.PlannedValue
		totalEV += (child.Progress / 100.0) * child.PlannedValue
		totalAC += child.ActualCost
		totalEstimatedEffort += child.EstimatedEffort
		totalActualEffort += child.ActualEffort
	}

	// Update Progress: Weighted by PlannedValue if PV exists, otherwise simple average
	if totalPV > 0 {
		parent.Progress = (totalEV / totalPV) * 100.0
	} else {
		// Fallback to simple average if no PV defined
		var sumProgress float64
		for _, child := range children {
			sumProgress += child.Progress
		}
		parent.Progress = sumProgress / float64(len(children))
	}

	parent.PlannedValue = totalPV
	parent.ActualCost = totalAC
	parent.EstimatedEffort = totalEstimatedEffort
	parent.ActualEffort = totalActualEffort

	// 4. Update the parent in DB
	err = s.repo.UpdateNode(ctx, parent)
	if err != nil {
		return err
	}

	// 5. Recursively roll up to the next parent
	parts := strings.Split(parentPath, ".")
	if len(parts) > 1 {
		grandParentPath := strings.Join(parts[:len(parts)-1], ".")
		return s.applyRollup(ctx, projectID, grandParentPath)
	}

	return nil
}
