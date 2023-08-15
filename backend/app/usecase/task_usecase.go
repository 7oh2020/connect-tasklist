package usecase

import (
	"context"
	"html"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
)

// タスクの操作
type ITaskUsecase interface {
	FindTasksByUserID(ctx context.Context, userID *dto.IDParam) ([]*entity.Task, error)
	CreateTask(ctx context.Context, arg *dto.CreateTaskParams) (string, error)
	ChangeTaskName(ctx context.Context, arg *dto.ChangeTaskNameParams) error
	CompleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error
	UncompleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error
	DeleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error
}

type TaskUsecase struct {
	service.ITaskService
}

func NewTaskUsecase(srv service.ITaskService) *TaskUsecase {
	return &TaskUsecase{srv}
}

func (u *TaskUsecase) FindTasksByUserID(ctx context.Context, userID *dto.IDParam) ([]*entity.Task, error) {
	if err := userID.Validate(); err != nil {
		return nil, err
	}
	return u.ITaskService.FindTasksByUserID(ctx, userID.Value())
}

func (u *TaskUsecase) CreateTask(ctx context.Context, arg *dto.CreateTaskParams) (string, error) {
	if err := arg.Validate(); err != nil {
		return "", err
	}
	return u.ITaskService.CreateTask(ctx, arg.UserID(), html.EscapeString(arg.Name()))
}

func (u *TaskUsecase) ChangeTaskName(ctx context.Context, arg *dto.ChangeTaskNameParams) error {
	if err := arg.Validate(); err != nil {
		return err
	}
	return u.ITaskService.ChangeTaskName(ctx, arg.ID(), arg.UserID(), html.EscapeString(arg.Name()))
}

func (u *TaskUsecase) CompleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error {
	if err := id.Validate(); err != nil {
		return err
	}
	if err := userID.Validate(); err != nil {
		return err
	}
	return u.ITaskService.CompleteTask(ctx, id.Value(), userID.Value())
}

func (u *TaskUsecase) UncompleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error {
	if err := id.Validate(); err != nil {
		return err
	}
	if err := userID.Validate(); err != nil {
		return err
	}
	return u.ITaskService.UncompleteTask(ctx, id.Value(), userID.Value())
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id *dto.IDParam, userID *dto.IDParam) error {
	if err := id.Validate(); err != nil {
		return err
	}
	if err := userID.Validate(); err != nil {
		return err
	}
	return u.ITaskService.DeleteTask(ctx, id.Value(), userID.Value())
}
