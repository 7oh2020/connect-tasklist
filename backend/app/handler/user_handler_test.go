package handler

import (
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1/userv1connect"
)

func TestUserHandler_NewUserHandler(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ userv1connect.UserServiceHandler = (*UserHandler)(nil)
	})
}
