package repository

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
)

// UserEntityの永続化を行う
type IUserRepository interface {
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
