package dto

import (
	"regexp"

	"github.com/7oh2020/connect-tasklist/backend/app"
)

type LoginParams struct {
	email    string
	password string
}

func NewLoginParams(email string, password string) *LoginParams {
	return &LoginParams{email, password}
}

func (f *LoginParams) Email() string {
	return f.email
}

func (f *LoginParams) Password() string {
	return f.password
}

func (f *LoginParams) Validate() error {
	if len([]rune(f.email)) > 100 {
		return &app.ErrInputValidationFailed{Msg: "email must be 100 characters or less"}
	}
	if len([]rune(f.password)) > 100 {
		return &app.ErrInputValidationFailed{Msg: "password must be 100 characters or less"}
	}
	// emailのフォーマットを検証する
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if ok := emailRegexp.MatchString(f.email); !ok {
		return &app.ErrInputValidationFailed{Msg: "invalid email"}
	}
	return nil
}
