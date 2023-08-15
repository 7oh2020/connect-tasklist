package service

import (
	"context"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestUserService_NewUserService(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IUserService = (*UserService)(nil)
	})
}

func TestUserService_FindUserByID(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := "id"
	user := &entity.User{
		ID:        value.NewID(id),
		Email:     value.NewEmail("email"),
		Password:  value.NewPassword("pass"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByID", ctx, id).Return(user, nil)
		srv := NewUserService(repo)
		ret, err := srv.FindUserByID(ctx, id)

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret.ID.Value())
		require.Equal(t, user.Email.Value(), ret.Email.Value())
		require.Equal(t, user.Password.Value(), ret.Password.Value())
		require.Equal(t, user.CreatedAt, ret.CreatedAt)
		require.Equal(t, user.UpdatedAt, ret.UpdatedAt)
		repo.AssertExpectations(t)
	})
	tt.Run("準正常系: UserIDが空の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		id := ""
		repo := new(mocks.IUserRepository)
		srv := NewUserService(repo)
		_, err := srv.FindUserByID(ctx, id)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないUserIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "user not found"}
		id := "another"
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByID", ctx, id).Return(nil, errExp)
		srv := NewUserService(repo)
		_, err := srv.FindUserByID(ctx, id)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
	})
}
