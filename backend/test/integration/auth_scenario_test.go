package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/interfaces/di"
	auth_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1/auth_v1connect"
	"github.com/stretchr/testify/require"
)

func TestAuthScenario(t *testing.T) {
	// テストサーバーの起動
	authHdr, err := di.InitAuth(issuer, keyPath, qry, timeout)
	require.NoError(t, err, "エラーが発生しないこと")
	mux := http.NewServeMux()
	mux.Handle(auth_v1connect.NewAuthServiceHandler(authHdr))
	ts := newTestServer(t, mux)
	defer ts.Close()

	email := "test@example.com"
	pass := "pass"

	// Login: Emailが空の場合
	res, err := ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, "", pass))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// Login: Emailが長すぎる場合
	res, err = ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, strings.Repeat("*", 101), pass))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// Login: Passwordが空の場合
	res, err = ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 401, res.status, "認証エラーになること")

	// Login: Passwordが長すぎる場合
	res, err = ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, strings.Repeat("*", 101)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// Login: Passwordが違う場合
	res, err = ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, "another"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 401, res.status, "認証エラーになること")

	// Login: 正しい入力の場合
	res, err = ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, pass))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// トークンを取得する
	var data auth_v1.LoginResponse
	err = json.Unmarshal([]byte(res.body), &data)
	require.NoError(t, err, "エラーが発生しないこと")
	require.NotEmpty(t, data.Token, "トークンが取得できること")
}
