package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	user_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1/user_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_NewUserHandler(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ user_v1connect.UserServiceHandler = (*UserHandler)(nil)
	})
}

func TestUserHandler_GetUser(tt *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	id := "id"
	user := &entity.User{
		ID:        value.NewID(id),
		Email:     value.NewEmail("test@example.com"),
		Password:  value.NewPassword("pass"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	arg := &user_v1.GetUserRequest{UserId: id}
	param := dto.NewIDParam(arg.UserId)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: ユーザーが存在しない場合", &domain.ErrNotFound{}, "not_found"},
		{"準正常系: その他のエラーの場合", &domain.ErrQueryFailed{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.IUserUsecase)
			if v.err == nil {
				uc.On("FindUserByID", ctx, param).Return(user, nil)
			} else {
				uc.On("FindUserByID", ctx, param).Return(nil, v.err)
			}
			hdr := NewUserHandler(uc)
			ret, err := hdr.GetUser(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
				require.Equal(t, id, ret.Msg.User.UserId)
				require.Equal(t, user.Email.Value(), ret.Msg.User.Email)
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
		})
	}
}
