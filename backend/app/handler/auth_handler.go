package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	auth_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1"
)

// AuthServiceHandlerの実装
type AuthHandler struct {
	usecase.IAuthUsecase
}

func NewAuthHandler(uc usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{uc}
}

func (h *AuthHandler) Login(ctx context.Context, arg *connect.Request[auth_v1.LoginRequest]) (*connect.Response[auth_v1.LoginResponse], error) {
	params := dto.NewLoginParams(arg.Msg.Email, arg.Msg.Password)
	info, err := h.IAuthUsecase.Login(ctx, params)
	if err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *app.ErrLoginFailed:
			return nil, connect.NewError(connect.CodeUnauthenticated, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeUnauthenticated, e)
		case *app.ErrInternal:
			return nil, connect.NewError(connect.CodeInternal, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&auth_v1.LoginResponse{
		Token: info.Token(),
	}), nil
}
