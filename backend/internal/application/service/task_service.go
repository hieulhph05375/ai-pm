package service

import (
	"context"
	"fmt"
	"project-mgmt/backend/internal/domain/entity"
	"project-mgmt/backend/internal/domain/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *entity.Task, actorID uint) error
	GetTask(ctx context.Context, id uint) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task, actorID uint) error
	DeleteTask(ctx context.Context, id uint) error
	ListTasksByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Task, int, error)
	ListActivities(ctx context.Context, taskID uint) ([]entity.TaskActivity, error)
	AddComment(ctx context.Context, taskID uint, actorID uint, content string) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(tr repository.TaskRepository) TaskService {
	return &taskService{taskRepo: tr}
}

func (s *taskService) CreateTask(ctx context.Context, task *entity.Task, actorID uint) error {
	if task.Status == "" {
		task.Status = entity.TaskStatusTodo
	}
	if task.Priority == "" {
		task.Priority = entity.TaskPriorityMedium
	}
	if task.Labels == nil {
		task.Labels = []string{}
	}
	task.CreatedBy = &actorID
	if task.AssigneeID == nil {
		task.AssigneeID = &actorID
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return err
	}

	return s.taskRepo.LogActivity(ctx, &entity.TaskActivity{
		TaskID:   task.ID,
		ActorID:  &actorID,
		Action:   "created",
		NewValue: task.Title,
	})
}

func (s *taskService) GetTask(ctx context.Context, id uint) (*entity.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *taskService) UpdateTask(ctx context.Context, task *entity.Task, actorID uint) error {
	old, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return err
	}

	// Log status change if any
	if old.Status != task.Status {
		_ = s.taskRepo.LogActivity(ctx, &entity.TaskActivity{
			TaskID:   task.ID,
			ActorID:  &actorID,
			Action:   "status_changed",
			OldValue: old.Status,
			NewValue: task.Status,
		})
	} else {
		_ = s.taskRepo.LogActivity(ctx, &entity.TaskActivity{
			TaskID:   task.ID,
			ActorID:  &actorID,
			Action:   "updated",
			NewValue: fmt.Sprintf("title: %s", task.Title),
		})
	}

	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, id uint) error {
	return s.taskRepo.Delete(ctx, id)
}

func (s *taskService) ListTasksByUser(ctx context.Context, userID uint, limit, offset int) ([]entity.Task, int, error) {
	tasks, err := s.taskRepo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.taskRepo.CountByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (s *taskService) ListActivities(ctx context.Context, taskID uint) ([]entity.TaskActivity, error) {
	return s.taskRepo.ListActivities(ctx, taskID)
}

func (s *taskService) AddComment(ctx context.Context, taskID uint, actorID uint, content string) error {
	return s.taskRepo.LogActivity(ctx, &entity.TaskActivity{
		TaskID:   taskID,
		ActorID:  &actorID,
		Action:   "COMMENT",
		NewValue: content,
	})
}
