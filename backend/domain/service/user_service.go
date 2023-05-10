package service

import (
	"context"
	"errors"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
)

// ユーザーのドメインロジック
type IUserService interface {
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
}

type UserService struct {
	repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("error: invalid parameter")
	}
	return s.IUserRepository.FindUserByID(ctx, id)
}
