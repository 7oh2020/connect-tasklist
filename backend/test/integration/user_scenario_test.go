package integration

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/interfaces/di"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/user/v1/user_v1connect"
	"github.com/stretchr/testify/require"
)

func TestUserScenario(t *testing.T) {
	// テストサーバーの起動
	userHdr := di.InitUser(qry)
	mux := http.NewServeMux()
	mux.Handle(user_v1connect.NewUserServiceHandler(userHdr))
	ts := newTestServer(t, mux)
	defer ts.Close()

	uid := "test"

	// GetUser: UserIDが空の場合
	res, err := ts.sendPostRequest(t, "", "/rpc.user.v1.UserService/GetUser", fmt.Sprintf(`{"user_id":"%s"}`, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// GetUser: UserIDが長すぎる場合
	res, err = ts.sendPostRequest(t, "", "/rpc.user.v1.UserService/GetUser", fmt.Sprintf(`{"user_id":"%s"}`, strings.Repeat("*", 51)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// GetUser: 正しい入力の場合
	res, err = ts.sendPostRequest(t, "", "/rpc.user.v1.UserService/GetUser", fmt.Sprintf(`{"user_id":"%s"}`, uid))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")
}
