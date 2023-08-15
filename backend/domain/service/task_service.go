package service

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
	"github.com/7oh2020/connect-tasklist/backend/util/clock"
	"github.com/7oh2020/connect-tasklist/backend/util/identification"
)

// タスクのドメインロジック
type ITaskService interface {
	FindTaskByID(ctx context.Context, id string) (*entity.Task, error)
	FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error)
	CreateTask(ctx context.Context, userID string, name string) (string, error)
	ChangeTaskName(ctx context.Context, id string, userID string, name string) error
	CompleteTask(ctx context.Context, id string, userID string) error
	UncompleteTask(ctx context.Context, id string, userID string) error
	DeleteTask(ctx context.Context, id, userID string) error
}

type TaskService struct {
	repository.ITaskRepository
	identification.IIDManager
	clock.IClockManager
}

func NewTaskService(repo repository.ITaskRepository, idManager identification.IIDManager, clockManager clock.IClockManager) *TaskService {
	return &TaskService{repo, idManager, clockManager}
}

func (s *TaskService) FindTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	if err := value.NewID(id).Validate(); err != nil {
		return nil, err
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return nil, &domain.ErrNotFound{Msg: "task not found"}
	}
	return task, nil
}

func (s *TaskService) FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error) {
	if err := value.NewID(userID).Validate(); err != nil {
		return nil, err
	}
	tasks, err := s.ITaskRepository.FindTasksByUserID(ctx, userID)
	if err != nil {
		return nil, &domain.ErrQueryFailed{}
	}
	return tasks, nil
}

func (s *TaskService) CreateTask(ctx context.Context, userID string, name string) (string, error) {
	now := s.IClockManager.GetNow()
	arg := &entity.Task{
		ID:          value.NewID(s.IIDManager.GenerateID()),
		UserID:      value.NewID(userID),
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := arg.Validate(); err != nil {
		return "", err
	}
	createdID, err := s.ITaskRepository.CreateTask(ctx, arg)
	if err != nil {
		return "", &domain.ErrQueryFailed{}
	}
	return createdID, nil
}

func (s *TaskService) ChangeTaskName(ctx context.Context, id string, userID string, name string) error {
	if err := value.NewID(id).Validate(); err != nil {
		return err
	}
	if err := value.NewID(userID).Validate(); err != nil {
		return err
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return &domain.ErrNotFound{Msg: "task not found"}
	}
	if !task.UserID.Equal(userID) {
		return &domain.ErrPermissionDenied{}
	}
	task.Name = name
	task.UpdatedAt = s.IClockManager.GetNow()
	if err := task.Validate(); err != nil {
		return err
	}
	if err := s.ITaskRepository.UpdateTask(ctx, task); err != nil {
		return &domain.ErrQueryFailed{}
	}
	return nil
}

func (s *TaskService) CompleteTask(ctx context.Context, id string, userID string) error {
	if err := value.NewID(id).Validate(); err != nil {
		return err
	}
	if err := value.NewID(userID).Validate(); err != nil {
		return err
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return &domain.ErrNotFound{Msg: "task not found"}
	}
	if !task.UserID.Equal(userID) {
		return &domain.ErrPermissionDenied{}
	}
	task.IsCompleted = true
	task.UpdatedAt = s.IClockManager.GetNow()
	if err := task.Validate(); err != nil {
		return err
	}
	if err := s.ITaskRepository.UpdateTask(ctx, task); err != nil {
		return &domain.ErrQueryFailed{}
	}
	return nil
}

func (s *TaskService) UncompleteTask(ctx context.Context, id string, userID string) error {
	if err := value.NewID(id).Validate(); err != nil {
		return err
	}
	if err := value.NewID(userID).Validate(); err != nil {
		return err
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return &domain.ErrNotFound{Msg: "task not found"}
	}
	if !task.UserID.Equal(userID) {
		return &domain.ErrPermissionDenied{}
	}
	task.IsCompleted = false
	task.UpdatedAt = s.IClockManager.GetNow()
	if err := task.Validate(); err != nil {
		return err
	}
	if err := s.ITaskRepository.UpdateTask(ctx, task); err != nil {
		return &domain.ErrQueryFailed{}
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string, userID string) error {
	if err := value.NewID(id).Validate(); err != nil {
		return err
	}
	if err := value.NewID(userID).Validate(); err != nil {
		return err
	}
	task, err := s.ITaskRepository.FindTaskByID(ctx, id)
	if err != nil {
		return &domain.ErrNotFound{Msg: "task not found"}
	}
	if !task.UserID.Equal(userID) {
		return &domain.ErrPermissionDenied{}
	}
	if err := s.ITaskRepository.DeleteTask(ctx, id); err != nil {
		return &domain.ErrQueryFailed{}
	}
	return nil
}
