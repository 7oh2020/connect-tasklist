package repository

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
)

// TaskEntityの永続化を行う
type ITaskRepository interface {
	FindTaskByID(ctx context.Context, id string) (*entity.Task, error)
	FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error)
	CreateTask(ctx context.Context, arg *entity.Task) (string, error)
	UpdateTask(ctx context.Context, arg *entity.Task) error
	DeleteTask(ctx context.Context, id string) error
}
