package usecase

import (
	"context"
	"errors"
	"html"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
)

// タスクの操作
type ITaskUsecase interface {
	GetTaskList(ctx context.Context, userID string) ([]*entity.Task, error)
	CreateTask(ctx context.Context, id string, userID string, name string, now time.Time) (string, error)
	ChangeTaskName(ctx context.Context, id string, userID string, name string, now time.Time) error
	CompleteTask(ctx context.Context, id string, userID string, now time.Time) error
	UncompleteTask(ctx context.Context, id string, userID string, now time.Time) error
	DeleteTask(ctx context.Context, id string, userID string) error
}

type TaskUsecase struct {
	service.ITaskService
}

func NewTaskUsecase(srv service.ITaskService) *TaskUsecase {
	return &TaskUsecase{srv}
}

func (u *TaskUsecase) GetTaskList(ctx context.Context, userID string) ([]*entity.Task, error) {
	if userID == "" || len([]rune(userID)) > 50 {
		return nil, errors.New("error: invalid request parameter")
	}
	return u.ITaskService.FindTasksByUserID(ctx, userID)
}

func (u *TaskUsecase) CreateTask(ctx context.Context, id string, userID string, name string, now time.Time) (string, error) {
	if id == "" || len([]rune(id)) > 50 || userID == "" || len([]rune(userID)) > 50 || name == "" || len([]rune(name)) > 100 {
		return "", errors.New("error: invalid request parameter")
	}
	// HTML文字列をエスケープする
	name = html.EscapeString(name)

	return u.ITaskService.CreateTask(ctx, id, userID, name, now)
}

func (u *TaskUsecase) ChangeTaskName(ctx context.Context, id string, userID string, name string, now time.Time) error {
	if id == "" || len([]rune(id)) > 50 || userID == "" || len([]rune(userID)) > 50 || name == "" || len([]rune(name)) > 100 {
		return errors.New("error: invalid request parameter")
	}
	// HTML文字列をエスケープする
	name = html.EscapeString(name)

	return u.ITaskService.ChangeTaskName(ctx, id, userID, name, now)
}

func (u *TaskUsecase) CompleteTask(ctx context.Context, id string, userID string, now time.Time) error {
	if id == "" || len([]rune(id)) > 50 || userID == "" || len([]rune(userID)) > 50 {
		return errors.New("error: invalid request parameter")
	}
	return u.ITaskService.CompleteTask(ctx, id, userID, now)
}

func (u *TaskUsecase) UncompleteTask(ctx context.Context, id string, userID string, now time.Time) error {
	if id == "" || len([]rune(id)) > 50 || userID == "" || len([]rune(userID)) > 50 {
		return errors.New("error: invalid request parameter")
	}
	return u.ITaskService.UncompleteTask(ctx, id, userID, now)
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id string, userID string) error {
	if id == "" || len([]rune(id)) > 50 || userID == "" || len([]rune(userID)) > 50 {
		return errors.New("error: invalid request parameter")
	}
	return u.ITaskService.DeleteTask(ctx, id, userID)
}
