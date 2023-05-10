package usecase

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app/util/auth"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

// ユーザーの認証処理
type IAuthUsecase interface {
	Login(ctx context.Context, email string, password string) (string, *entity.User, error)
}

type AuthUsecase struct {
	repository.IUserRepository
	auth.ITokenManager
	duration time.Duration
}

func NewAuthUsecase(repo repository.IUserRepository, tm auth.ITokenManager, duration time.Duration) *AuthUsecase {
	return &AuthUsecase{repo, tm, duration}
}

func (u *AuthUsecase) Login(ctx context.Context, email string, password string) (string, *entity.User, error) {
	// 入力をバリデーションする
	if email == "" || len([]rune(email)) > 100 || password == "" || len([]rune(password)) > 100 {
		return "", nil, errors.New("error: invalid parameter")
	}
	// emailのフォーマットをバリデーションする
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if ok := emailRegexp.MatchString(email); !ok {
		return "", nil, errors.New("error: invalid email")
	}
	res, err := u.IUserRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	// bcrypt方式でパスワードが一致するか検証する
	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(password)); err != nil {
		return "", nil, errors.New("error: password does not match")
	}
	// JWTを作成する
	token, err := u.ITokenManager.CreateToken(res.ID, u.duration)
	if err != nil {
		return "", nil, err
	}
	return token, res, nil
}
