package handler

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	auth_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1/auth_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_NewAuthHandler(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ auth_v1connect.AuthServiceHandler = (*AuthHandler)(nil)
	})
}

func TestAuthHandler_Login(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	email := "test@example.com"
	pass := "pass"
	token := "token"
	info := dto.NewUserInfo(id, email, token)
	arg := &auth_v1.LoginRequest{Email: email, Password: pass}
	params := dto.NewLoginParams(arg.Email, arg.Password)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		ret     *dto.UserInfo
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", info, nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", nil, &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", nil, &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ユーザーが存在しない場合", nil, &domain.ErrNotFound{}, "unauthenticated"},
		{"準正常系: 認証に失敗した場合", nil, &app.ErrLoginFailed{}, "unauthenticated"},
		{"準正常系: アプリ内部エラーの場合", nil, &app.ErrInternal{}, "internal"},
		{"準正常系: その他のエラーの場合", nil, &domain.ErrQueryFailed{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.IAuthUsecase)
			if v.err == nil {
				uc.On("Login", ctx, params).Return(info, nil)
			} else {
				uc.On("Login", ctx, params).Return(nil, v.err)
			}
			hdr := NewAuthHandler(uc)
			ret, err := hdr.Login(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
				require.Equal(t, token, ret.Msg.Token)
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
		})
	}
}
