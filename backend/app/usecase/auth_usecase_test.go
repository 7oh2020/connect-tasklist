package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestAuthUsecase_NewAuthUsecase(tt *testing.T) {
	tt.Run("正常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IAuthUsecase = (*AuthUsecase)(nil)
	})
}

func TestAuthUsecase_Login(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	email := "admin@example.com"
	pass := "pass"
	duration := 1 * time.Hour
	token := "token"
	user := &entity.User{
		ID:        "id",
		Email:     email,
		Password:  "$2a$10$YfHxWNfL8Ba2ltl6TRHMVuN0WPXxAuB5L7w1Y0jqaFcn2bUDoUq9W",
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo := new(mocks.IUserRepository)
	repo.On("FindUserByEmail", ctx, email).Return(user, nil)

	tm := new(mocks.ITokenManager)
	tm.On("CreateToken", user.ID, duration).Return(token, nil)

	testcases := []struct {
		title    string
		email    string
		pass     string
		resToken string
		resUser  *entity.User
		err      error
	}{
		{title: "正常系: 正しい入力の場合", email: email, pass: pass, resToken: token, resUser: user, err: nil},
		{title: "準正常系: emailが空の場合", email: "", pass: pass, resToken: token, resUser: nil, err: errors.New("error: invalid parameter")},
		{title: "準正常系: passwordが空の場合", email: email, pass: "", resToken: token, resUser: nil, err: errors.New("error: invalid parameter")},
		{title: "準正常系: emailのフォーマットが不正の場合", email: "email", pass: pass, resToken: token, resUser: nil, err: errors.New("error: invalid email")},
		{title: "準正常系: パスワードが一致しない場合", email: email, pass: "another", resToken: token, resUser: nil, err: errors.New("error: password does not match")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewAuthUsecase(repo, tm, duration)
			tok, usr, err := uc.Login(ctx, tc.email, tc.pass)
			require.Equal(t, tc.err, err)

			if err == nil {
				repo.AssertExpectations(t)
				require.Equal(t, tc.resToken, tok)
				require.Equal(t, tc.resUser.ID, usr.ID)
				require.Equal(t, tc.resUser.Email, usr.Email)
				require.Equal(t, tc.resUser.Password, usr.Password)
				require.Equal(t, tc.resUser.CreatedAt, usr.CreatedAt)
				require.Equal(t, tc.resUser.UpdatedAt, usr.UpdatedAt)
			}
		})
	}

}
