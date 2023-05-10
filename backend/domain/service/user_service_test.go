package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestUserService_NewUserService(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IUserService = (*UserService)(nil)
	})
}

func TestUserService_FindUserByID(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := "id"
	user := &entity.User{
		ID:        id,
		Email:     "email",
		Password:  "pass",
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo := new(mocks.IUserRepository)
	repo.On("FindUserByID", ctx, id).Return(user, nil)

	testcases := []struct {
		title string
		id    string
		res   *entity.User
		err   error
	}{
		{title: "正常系: 正しい入力の場合", id: id, res: user, err: nil},
		{title: "準正常系: idが空の場合", id: "", res: user, err: errors.New("error: invalid parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewUserService(repo)
			res, err := srv.FindUserByID(ctx, tc.id)
			require.Equal(t, tc.err, err)

			if res != nil {
				repo.AssertExpectations(t)
				require.Equal(t, tc.res.ID, res.ID)
				require.Equal(t, tc.res.Email, res.Email)
				require.Equal(t, tc.res.Password, res.Password)
				require.Equal(t, tc.res.CreatedAt, res.CreatedAt)
				require.Equal(t, tc.res.UpdatedAt, res.UpdatedAt)
			}
		})
	}
}
