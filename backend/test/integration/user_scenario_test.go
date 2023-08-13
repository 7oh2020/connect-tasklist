package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/interfaces/di"
	userv1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1/userv1connect"
	"github.com/stretchr/testify/require"
)

func TestUserScenario_Login(t *testing.T) {
	// テストサーバーの起動
	mux := http.NewServeMux()
	userHdr, err := di.InitUser(issuer, keyPath, qry, timeout)
	require.NoError(t, err, "エラーが発生しないこと")
	mux.Handle(userv1connect.NewUserServiceHandler(userHdr))
	ts := newTestServer(t, mux)
	defer ts.Close()

	email := "test@example.com"
	pass := "pass"

	// Emailが空の場合
	res, err := ts.sendPostRequest(t, "/rpc.user.v1.UserService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, "", pass))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 401, res.status, "認証エラーになること")

	// Passwordが空の場合
	res, err = ts.sendPostRequest(t, "/rpc.user.v1.UserService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 401, res.status, "認証エラーになること")

	// 正しい入力の場合
	res, err = ts.sendPostRequest(t, "/rpc.user.v1.UserService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, pass))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	var data userv1.LoginResponse
	err = json.Unmarshal([]byte(res.body), &data)
	require.NoError(t, err, "エラーが発生しないこと")
	require.NotEmpty(t, data.Token, "トークンが取得できること")
}
