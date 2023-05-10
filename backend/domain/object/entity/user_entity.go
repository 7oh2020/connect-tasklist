package entity

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id string, email string, password string, now time.Time) (*User, error) {
	user := &User{
		ID:        id,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}

// 自身のフィールドをバリデーションする
func (u *User) Validate() error {
	if u.ID == "" || u.Email == "" || u.Password == "" {
		return errors.New("error: validation failed")
	}
	return nil
}
