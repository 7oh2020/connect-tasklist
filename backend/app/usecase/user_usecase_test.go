package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_NewUserUsecase(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装して	いるか", func(t *testing.T) {
		var _ IUserUsecase = (*UserUsecase)(nil)
	})
}

func TestUserUsecase_FindUserByID(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	now := time.Now().UTC()
	user := &entity.User{
		ID:        value.NewID(id),
		Email:     value.NewEmail("test@example.com"),
		Password:  value.NewPassword("pass"),
		CreatedAt: now,
		UpdatedAt: now,
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		srv := new(mocks.IUserService)
		srv.On("FindUserByID", ctx, id).Return(user, nil)
		uc := NewUserUsecase(srv)
		ret, err := uc.FindUserByID(ctx, dto.NewIDParam(id))

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret.ID.Value())
		require.Equal(t, user.Email.Value(), ret.Email.Value())
		require.Equal(t, user.Password.Value(), ret.Password.Value())
		require.Equal(t, user.CreatedAt, ret.CreatedAt)
		require.Equal(t, user.UpdatedAt, ret.UpdatedAt)
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
		id := strings.Repeat("*", 51)
		srv := new(mocks.IUserService)
		uc := NewUserUsecase(srv)
		_, err := uc.FindUserByID(ctx, dto.NewIDParam(id))

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)
	})
}
