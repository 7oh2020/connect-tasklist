package usecase

import (
	"context"
	"errors"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
)

// ユーザーの操作
type IUserUsecase interface {
	GetUser(ctx context.Context, id string) (*entity.User, error)
}

type UserUsecase struct {
	service.IUserService
}

func NewUserUsecase(srv service.IUserService) *UserUsecase {
	return &UserUsecase{srv}
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	if id == "" || len([]rune(id)) > 50 {
		return nil, errors.New("error: invalid request parameter")
	}
	return u.IUserService.FindUserByID(ctx, id)
}
