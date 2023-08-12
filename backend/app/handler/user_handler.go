package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	userv1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1"
)

// UserServiceHandlerの実装
type UserHandler struct {
	usecase.IAuthUsecase
	usecase.IUserUsecase
}

func NewUserHandler(uca usecase.IAuthUsecase, ucu usecase.IUserUsecase) *UserHandler {
	return &UserHandler{uca, ucu}
}

func (h *UserHandler) GetUser(ctx context.Context, arg *connect.Request[userv1.GetUserRequest]) (*connect.Response[userv1.GetUserResponse], error) {
	user, err := h.IUserUsecase.GetUser(ctx, arg.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&userv1.GetUserResponse{
		User: &userv1.User{
			UserId: user.ID,
			Email:  user.Email,
		},
	}), nil
}

func (h *UserHandler) Login(ctx context.Context, arg *connect.Request[userv1.LoginRequest]) (*connect.Response[userv1.LoginResponse], error) {
	token, user, err := h.IAuthUsecase.Login(ctx, arg.Msg.Email, arg.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	return connect.NewResponse(&userv1.LoginResponse{
		User: &userv1.User{
			UserId: user.ID,
			Email:  user.Email,
		},
		Token: token,
	}), nil
}
