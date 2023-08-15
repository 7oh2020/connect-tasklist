package usecase

import (
	"context"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	"github.com/7oh2020/connect-tasklist/backend/util/auth"
	"golang.org/x/crypto/bcrypt"
)

// ユーザーの認証処理
type IAuthUsecase interface {
	Login(ctx context.Context, arg *dto.LoginParams) (*dto.UserInfo, error)
}

type AuthUsecase struct {
	repository.IUserRepository
	auth.ITokenManager
	timeout time.Duration
}

func NewAuthUsecase(repo repository.IUserRepository, tm auth.ITokenManager, timeout time.Duration) *AuthUsecase {
	return &AuthUsecase{repo, tm, timeout}
}

func (u *AuthUsecase) Login(ctx context.Context, arg *dto.LoginParams) (*dto.UserInfo, error) {
	if err := arg.Validate(); err != nil {
		return nil, err
	}
	user, err := u.IUserRepository.FindUserByEmail(ctx, arg.Email())
	if err != nil {
		return nil, &domain.ErrNotFound{Msg: "user not found"}
	}
	// bcrypt方式でパスワードが一致するか検証する
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.Value()), []byte(arg.Password())); err != nil {
		return nil, &app.ErrLoginFailed{Msg: "password does not match"}
	}
	// JWTを作成する
	token, err := u.ITokenManager.CreateToken(user.ID.Value(), u.timeout)
	if err != nil {
		return nil, &app.ErrInternal{Msg: "failed to create token"}
	}
	return dto.NewUserInfo(user.ID.Value(), user.Email.Value(), token), nil
}
