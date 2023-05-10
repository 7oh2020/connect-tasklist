package service

import (
	"context"
	"errors"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
)

// タスクのドメインロジック
type ITaskService interface {
	FindTaskByID(ctx context.Context, id string) (*entity.Task, error)
	FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error)
	CreateTask(ctx context.Context, id string, userID string, name string, now time.Time) (string, error)
	ChangeTaskName(ctx context.Context, id string, userID string, name string, now time.Time) error
	CompleteTask(ctx context.Context, id string, userID string, now time.Time) error
	UncompleteTask(ctx context.Context, id string, userID string, now time.Time) error
	DeleteTask(ctx context.Context, id, userID string) error
}

type TaskService struct {
	repository.ITaskRepository
}

func NewTaskService(repo repository.ITaskRepository) *TaskService {
	return &TaskService{repo}
}

func (s *TaskService) FindTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	if id == "" {
		return nil, errors.New("error: invalid parameter")
	}
	return s.ITaskRepository.FindTaskByID(ctx, id)
}

func (s *TaskService) FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error) {
	if userID == "" {
		return nil, errors.New("error: invalid parameter")
	}
	return s.ITaskRepository.FindTasksByUserID(ctx, userID)
}

func (s *TaskService) CreateTask(ctx context.Context, id string, userID string, name string, now time.Time) (string, error) {
	task, err := entity.NewTask(id, userID, name, now)
	if err != nil {
		return "", err
	}
	return s.ITaskRepository.CreateTask(ctx, task)
}

func (s *TaskService) ChangeTaskName(ctx context.Context, id string, userID string, name string, now time.Time) error {
	if id == "" || userID == "" || name == "" {
		return errors.New("error: invalid parameter")
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return errors.New("error: permission denied to update task")
	}
	task.Name = name
	task.UpdatedAt = now
	if err := task.Validate(); err != nil {
		return err
	}
	return s.ITaskRepository.UpdateTask(ctx, task)
}

func (s *TaskService) CompleteTask(ctx context.Context, id string, userID string, now time.Time) error {
	if id == "" || userID == "" {
		return errors.New("error: invalid parameter")
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return errors.New("error: permission denied to update task")
	}
	task.IsCompleted = true
	task.UpdatedAt = now
	if err := task.Validate(); err != nil {
		return err
	}
	return s.ITaskRepository.UpdateTask(ctx, task)
}

func (s *TaskService) UncompleteTask(ctx context.Context, id string, userID string, now time.Time) error {
	if id == "" || userID == "" {
		return errors.New("error: invalid parameter")
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return err
	}
	if task.UserID != userID {
		return errors.New("error: permission denied to update task")
	}
	task.IsCompleted = false
	task.UpdatedAt = now
	if err := task.Validate(); err != nil {
		return err
	}
	return s.ITaskRepository.UpdateTask(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id string, userID string) error {
	if id == "" || userID == "" {
		return errors.New("error: invalid parameter")
	}
	if task, err := s.ITaskRepository.FindTaskByID(ctx, id); err != nil {
		return err
	} else if task.UserID != userID {
		return errors.New("error: permission denied to delete task")
	}
	return s.ITaskRepository.DeleteTask(ctx, id)
}
