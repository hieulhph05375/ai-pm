package service

import (
	"context"
	"errors"
	"fmt"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type TimesheetService interface {
	CreateTimesheet(ctx context.Context, timesheet *entity.Timesheet) error
	GetTimesheet(ctx context.Context, id int) (*entity.Timesheet, error)
	UpdateTimesheet(ctx context.Context, timesheet *entity.Timesheet) error
	DeleteTimesheet(ctx context.Context, id int) error
	ListByUser(ctx context.Context, userID int, limit, offset int) ([]entity.Timesheet, int, error)
	ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Timesheet, int, error)
}

type timesheetService struct {
	tsRepo   repository.TimesheetRepository
	wbsRepo  repository.WBSRepository
	notifSvc NotificationService
}

func NewTimesheetService(tr repository.TimesheetRepository, wr repository.WBSRepository, ns NotificationService) TimesheetService {
	return &timesheetService{tsRepo: tr, wbsRepo: wr, notifSvc: ns}
}

func (s *timesheetService) CreateTimesheet(ctx context.Context, t *entity.Timesheet) error {
	if t.Hours <= 0 {
		return errors.New("hours must be greater than zero")
	}
	if t.Hours > 24 {
		return errors.New("hours cannot exceed 24 per day")
	}

	err := s.tsRepo.Create(ctx, t)
	if err != nil {
		return err
	}

	// Effort Roll-up calculation
	if t.NodeID != nil {
		node, err := s.wbsRepo.GetNodeByID(ctx, *t.NodeID)
		if err == nil && node != nil {
			node.ActualEffort += t.Hours
			_ = s.wbsRepo.UpdateNode(ctx, node)

			// Trigger EFFORT_OVERRUN notification if actual > estimated and estimated > 0
			if node.EstimatedEffort > 0 && node.ActualEffort > node.EstimatedEffort {
				nodeRefType := "wbs_node"
				_ = s.notifSvc.Create(ctx, t.UserID,
					"EFFORT_OVERRUN",
					"Vượt effort ước tính: "+node.Title,
					fmt.Sprintf("Task '%s' đã ghi %.1f giờ thực tế, vượt %.1f giờ ước tính.", node.Title, node.ActualEffort, node.EstimatedEffort),
					&node.ID, &nodeRefType,
				)
			}
		}
	}

	return nil
}

func (s *timesheetService) GetTimesheet(ctx context.Context, id int) (*entity.Timesheet, error) {
	return s.tsRepo.GetByID(ctx, id)
}

func (s *timesheetService) UpdateTimesheet(ctx context.Context, t *entity.Timesheet) error {
	if t.Hours <= 0 {
		return errors.New("hours must be greater than zero")
	}
	if t.Hours > 24 {
		return errors.New("hours cannot exceed 24 per day")
	}

	oldT, err := s.tsRepo.GetByID(ctx, t.ID)
	if err != nil {
		return err
	}

	err = s.tsRepo.Update(ctx, t)
	if err != nil {
		return err
	}

	// Effort Roll up Re-calculation
	// If a timesheet is updated, we withdraw old hours and inject new hours
	if oldT.NodeID != nil && t.NodeID != nil && *oldT.NodeID == *t.NodeID {
		node, err := s.wbsRepo.GetNodeByID(ctx, *t.NodeID)
		if err == nil && node != nil {
			node.ActualEffort -= oldT.Hours // remove old
			node.ActualEffort += t.Hours    // add new

			// Prevent negative effort just in case
			if node.ActualEffort < 0 {
				node.ActualEffort = 0
			}

			_ = s.wbsRepo.UpdateNode(ctx, node)
		}
	} else {
		// Handled moved tasks
		if oldT.NodeID != nil {
			oldNode, err := s.wbsRepo.GetNodeByID(ctx, *oldT.NodeID)
			if err == nil && oldNode != nil {
				oldNode.ActualEffort -= oldT.Hours
				if oldNode.ActualEffort < 0 {
					oldNode.ActualEffort = 0
				}
				_ = s.wbsRepo.UpdateNode(ctx, oldNode)
			}
		}
		if t.NodeID != nil {
			newNode, err := s.wbsRepo.GetNodeByID(ctx, *t.NodeID)
			if err == nil && newNode != nil {
				newNode.ActualEffort += t.Hours
				_ = s.wbsRepo.UpdateNode(ctx, newNode)
			}
		}
	}

	return nil
}

func (s *timesheetService) DeleteTimesheet(ctx context.Context, id int) error {
	oldT, err := s.tsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.tsRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Reverse Effort
	if oldT.NodeID != nil {
		node, err := s.wbsRepo.GetNodeByID(ctx, *oldT.NodeID)
		if err == nil && node != nil {
			node.ActualEffort -= oldT.Hours
			if node.ActualEffort < 0 {
				node.ActualEffort = 0
			}
			_ = s.wbsRepo.UpdateNode(ctx, node)
		}
	}

	return nil
}

func (s *timesheetService) ListByUser(ctx context.Context, userID int, limit, offset int) ([]entity.Timesheet, int, error) {
	return s.tsRepo.ListByUser(ctx, userID, limit, offset)
}

func (s *timesheetService) ListByProject(ctx context.Context, projectID int, limit, offset int) ([]entity.Timesheet, int, error) {
	return s.tsRepo.ListByProject(ctx, projectID, limit, offset)
}
