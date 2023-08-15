package usecase

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
)

// ユーザーの操作
type IUserUsecase interface {
	FindUserByID(ctx context.Context, id *dto.IDParam) (*entity.User, error)
}

type UserUsecase struct {
	service.IUserService
}

func NewUserUsecase(srv service.IUserService) *UserUsecase {
	return &UserUsecase{srv}
}

func (u *UserUsecase) FindUserByID(ctx context.Context, id *dto.IDParam) (*entity.User, error) {
	if err := id.Validate(); err != nil {
		return nil, err
	}
	return u.IUserService.FindUserByID(ctx, id.Value())
}
