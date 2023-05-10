package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_NewUserUsecase(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装して	いるか", func(t *testing.T) {
		var _ IUserUsecase = (*UserUsecase)(nil)
	})
}

func TestUserUsecase_GetUser(tt *testing.T) {
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	now := time.Now().UTC()
	user := &entity.User{
		ID:        id,
		Email:     "email",
		Password:  "pass",
		CreatedAt: now,
		UpdatedAt: now,
	}

	srv := new(mocks.IUserService)
	srv.On("FindUserByID", ctx, id).Return(user, nil)

	testcases := []struct {
		title string
		id    string
		res   *entity.User
		err   error
	}{
		{title: "正常系: 正しい入力の場合", id: id, res: user, err: nil},
		{title: "準正常系: idが空の場合", id: "", res: nil, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), res: nil, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), res: nil, err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewUserUsecase(srv)
			res, err := uc.GetUser(ctx, tc.id)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)

			if res != nil {
				srv.AssertExpectations(t)
				require.Equal(t, id, res.ID)
				assert.Equal(t, user.Email, res.Email)
				assert.Equal(t, user.Password, res.Password)
				assert.Equal(t, user.CreatedAt, res.CreatedAt)
				assert.Equal(t, user.UpdatedAt, res.UpdatedAt)
			}
		})
	}
}
