package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestAuthUsecase_NewAuthUsecase(tt *testing.T) {
	tt.Run("正常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IAuthUsecase = (*AuthUsecase)(nil)
	})
}

func TestAuthUsecase_Login(tt *testing.T) {
	ctx := context.Background()
	timeout := time.Hour
	id := "id"
	email := "test@example.com"
	pass := "pass"
	hPass := "$2a$10$YfHxWNfL8Ba2ltl6TRHMVuN0WPXxAuB5L7w1Y0jqaFcn2bUDoUq9W"
	now := time.Now().UTC()
	user := &entity.User{
		ID:        value.NewID(id),
		Email:     value.NewEmail(email),
		Password:  value.NewPassword(hPass),
		CreatedAt: now,
		UpdatedAt: now,
	}
	token := "token"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := dto.NewLoginParams(email, pass)
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByEmail", ctx, email).Return(user, nil)
		tm := new(mocks.ITokenManager)
		tm.On("CreateToken", id, timeout).Return(token, nil)
		uc := NewAuthUsecase(repo, tm, timeout)
		ret, err := uc.Login(ctx, arg)

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret.ID())
		require.Equal(t, email, ret.Email())
		require.Equal(t, token, ret.Token())
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "invalid email"}
		arg := dto.NewLoginParams("test", pass)
		repo := new(mocks.IUserRepository)
		tm := new(mocks.ITokenManager)
		uc := NewAuthUsecase(repo, tm, timeout)
		_, err := uc.Login(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないEmailの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "user not found"}
		arg := dto.NewLoginParams("another@example.com", pass)
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByEmail", ctx, arg.Email()).Return(nil, errExp)
		tm := new(mocks.ITokenManager)
		uc := NewAuthUsecase(repo, tm, timeout)
		_, err := uc.Login(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})
	tt.Run("準正常系: Passwordが一致しない場合", func(t *testing.T) {
		errExp := &app.ErrLoginFailed{Msg: "password does not match"}
		arg := dto.NewLoginParams(email, "another")
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByEmail", ctx, email).Return(user, nil)
		tm := new(mocks.ITokenManager)
		uc := NewAuthUsecase(repo, tm, timeout)
		_, err := uc.Login(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})
	tt.Run("準正常系: 内部エラーの場合", func(t *testing.T) {
		errExp := &app.ErrInternal{Msg: "failed to create token"}
		arg := dto.NewLoginParams(email, pass)
		repo := new(mocks.IUserRepository)
		repo.On("FindUserByEmail", ctx, email).Return(user, nil)
		tm := new(mocks.ITokenManager)
		tm.On("CreateToken", id, timeout).Return("", errExp)
		uc := NewAuthUsecase(repo, tm, timeout)
		_, err := uc.Login(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		tm.AssertExpectations(t)
	})
}
