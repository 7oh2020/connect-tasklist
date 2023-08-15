package service

import (
	"context"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
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
	if err := value.NewID(id).Validate(); err != nil {
		return nil, err
	}
	user, err := s.IUserRepository.FindUserByID(ctx, id)
	if err != nil {
		return nil, &domain.ErrNotFound{Msg: "user not found"}
	}
	return user, nil
}
