package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	user_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1"
)

// UserServiceHandlerの実装
type UserHandler struct {
	usecase.IUserUsecase
}

func NewUserHandler(uc usecase.IUserUsecase) *UserHandler {
	return &UserHandler{uc}
}

func (h *UserHandler) GetUser(ctx context.Context, arg *connect.Request[user_v1.GetUserRequest]) (*connect.Response[user_v1.GetUserResponse], error) {
	user, err := h.IUserUsecase.FindUserByID(ctx, dto.NewIDParam(arg.Msg.UserId))
	if err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeNotFound, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&user_v1.GetUserResponse{
		User: &user_v1.User{
			UserId: user.ID.Value(),
			Email:  user.Email.Value(),
		},
	}), nil
}
