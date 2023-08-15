package entity

import (
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
)

type User struct {
	ID        *value.ID
	Email     *value.Email
	Password  *value.Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// フィールドの妥当性を検証する
func (u *User) Validate() error {
	if err := u.ID.Validate(); err != nil {
		return err
	}
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if err := u.Password.Validate(); err != nil {
		return err
	}
	return nil
}
