package interceptor

import (
	"context"
	"errors"
	"strings"

	"github.com/7oh2020/connect-tasklist/backend/app/util/auth"
	"github.com/7oh2020/connect-tasklist/backend/app/util/contextkey"
	"github.com/bufbuild/connect-go"
)

// リクエストのJWTを検証する。成功時にはUserIDをコンテキストにセットする
func NewAuthInterceptor(issuer string, keyPath string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// リクエストヘッダーからJWTを取得する
			token := req.Header().Get("Authorization")
			if token == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("error: invalid token"))
			}
			token = strings.TrimPrefix(token, "Bearer")
			token = strings.TrimSpace(token)

			// トークンを検証しUserIDを取得する
			tm, err := auth.NewTokenManager(issuer, keyPath)
			if err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			uid, err := tm.GetUserID(token)
			if err != nil {
				return nil, connect.NewError(connect.CodeAborted, err)
			}

			// コンテキストにUserIDをセットする
			cw := contextkey.NewContextWriter()
			ctx = cw.SetUserID(ctx, uid)

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
