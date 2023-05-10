package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserEntity_NewUserEntity(tt *testing.T) {
	now := time.Now().UTC()
	id := "id"
	email := "email"
	pass := "pass"
	user := &User{
		ID:        id,
		Email:     email,
		Password:  pass,
		CreatedAt: now,
		UpdatedAt: now,
	}

	testcases := []struct {
		title string
		id    string
		email string
		pass  string
		res   *User
		err   error
	}{
		{title: "正常系: 正しい入力の場合", id: id, email: email, pass: pass, res: user, err: nil},
		{title: "準正常系: idが空の場合", id: "", email: email, pass: pass, res: nil, err: errors.New("error: validation failed")},
		{title: "準正常系: emailが空の場合", id: id, email: "", pass: pass, res: nil, err: errors.New("error: validation failed")},
		{title: "準正常系: passwordが空の場合", id: id, email: email, pass: "", res: nil, err: errors.New("error: validation failed")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			res, err := NewUser(tc.id, tc.email, tc.pass, now)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)
		})
	}

}
