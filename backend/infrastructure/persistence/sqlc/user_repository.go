package sqlc

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
)

// ユーザー永続化のSQLC実装
type SQLCUserRepository struct {
	db.Querier
}

func NewSQLCUserRepository(qry db.Querier) *SQLCUserRepository {
	return &SQLCUserRepository{qry}
}

func (r *SQLCUserRepository) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	res, err := r.Querier.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:        res.ID,
		Email:     res.Email,
		Password:  res.Password,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

func (r *SQLCUserRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	res, err := r.Querier.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:        res.ID,
		Email:     res.Email,
		Password:  res.Password,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
