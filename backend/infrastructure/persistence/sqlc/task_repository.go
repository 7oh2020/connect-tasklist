package sqlc

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
)

// タスク永続化のSQLC実装
type SQLCTaskRepository struct {
	db.Querier
}

func NewSQLCTaskRepository(qry db.Querier) *SQLCTaskRepository {
	return &SQLCTaskRepository{qry}
}

func (r *SQLCTaskRepository) FindTaskByID(ctx context.Context, id string) (*entity.Task, error) {
	res, err := r.Querier.FindTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Task{
		ID:          res.ID,
		UserID:      res.UserID,
		Name:        res.Name,
		IsCompleted: res.IsCompleted,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}

func (r *SQLCTaskRepository) FindTasksByUserID(ctx context.Context, userID string) ([]*entity.Task, error) {
	res, err := r.Querier.FindTasksByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	tasks := make([]*entity.Task, len(res))
	for i, v := range res {
		tasks[i] = &entity.Task{
			ID:          v.ID,
			UserID:      v.UserID,
			Name:        v.Name,
			IsCompleted: v.IsCompleted,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
	}
	return tasks, nil
}

func (r *SQLCTaskRepository) CreateTask(ctx context.Context, arg *entity.Task) (string, error) {
	return r.Querier.CreateTask(ctx, db.CreateTaskParams{
		ID:          arg.ID,
		UserID:      arg.UserID,
		Name:        arg.Name,
		IsCompleted: arg.IsCompleted,
		CreatedAt:   arg.CreatedAt,
		UpdatedAt:   arg.UpdatedAt,
	})
}

func (r *SQLCTaskRepository) UpdateTask(ctx context.Context, arg *entity.Task) error {
	return r.Querier.UpdateTask(ctx, db.UpdateTaskParams{
		ID:          arg.ID,
		Name:        arg.Name,
		IsCompleted: arg.IsCompleted,
		UpdatedAt:   arg.UpdatedAt,
	})
}

func (r *SQLCTaskRepository) DeleteTask(ctx context.Context, id string) error {
	return r.Querier.DeleteTask(ctx, id)
}
