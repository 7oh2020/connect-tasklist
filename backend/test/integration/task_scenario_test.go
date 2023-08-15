package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/di"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/interceptor"
	auth_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/auth/v1/auth_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1/task_v1connect"
	"github.com/stretchr/testify/require"
)

func TestTaskScenario(t *testing.T) {
	// テストサーバーの起動
	authInterceptor := connect.WithInterceptors(interceptor.NewAuthInterceptor(issuer, keyPath))
	taskHdr := di.InitTask(qry)
	authHdr, err := di.InitAuth(issuer, keyPath, qry, timeout)
	require.NoError(t, err, "エラーが発生しないこと")
	mux := http.NewServeMux()
	mux.Handle(auth_v1connect.NewAuthServiceHandler(authHdr))
	mux.Handle(task_v1connect.NewTaskServiceHandler(taskHdr, authInterceptor))
	ts := newTestServer(t, mux)
	defer ts.Close()

	anotherTaskID := "t3"

	// Login: ログインしてトークンを取得する
	res, err := ts.sendPostRequest(t, "", "/rpc.auth.v1.AuthService/Login", fmt.Sprintf(`{"email":"%s", "password":"%s"}`, "test@example.com", "pass"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// トークンを取得する
	var data auth_v1.LoginResponse
	err = json.Unmarshal([]byte(res.body), &data)
	require.NoError(t, err, "エラーが発生しないこと")
	token := data.Token

	// GetTaskList: トークンが不正な場合
	res, err = ts.sendPostRequest(t, "another", "/rpc.task.v1.TaskService/GetTaskList", "{}")
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 401, res.status, "認証エラーになること")

	// GetTaskList: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/GetTaskList", "{}")
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// CreateTask: Nameが空の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CreateTask", fmt.Sprintf(`{"name":"%s"}`, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// CreateTask: Nameが長すぎる場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CreateTask", fmt.Sprintf(`{"name":"%s"}`, strings.Repeat("*", 101)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// CreateTask: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CreateTask", fmt.Sprintf(`{"name":"%s"}`, "new task"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// 作成されたTaskIDが取得できること
	var task struct {
		CreatedID string `json:"createdId"`
	}
	err = json.Unmarshal([]byte(res.body), &task)
	require.NoError(t, err, "エラーが発生しないこと")
	taskID := task.CreatedID
	require.NotEmpty(t, taskID, "TaskIDが取得できること")

	// ChangeTaskName: Nameが空の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/ChangeTaskName", fmt.Sprintf(`{"task_id":"%s", "name":"%s"}`, taskID, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// ChangeTaskName: Nameが長すぎる場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/ChangeTaskName", fmt.Sprintf(`{"task_id":"%s", "name":"%s"}`, taskID, strings.Repeat("*", 101)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// ChangeTaskName: 存在しないTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/ChangeTaskName", fmt.Sprintf(`{"task_id":"%s", "name":"%s"}`, "another", "Update Task Name"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 404, res.status, "404エラーになること")

	// ChangeTaskName: 他人のTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/ChangeTaskName", fmt.Sprintf(`{"task_id":"%s", "name":"%s"}`, anotherTaskID, "Update Task Name"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 403, res.status, "パーミッションエラーになること")

	// ChangeTaskName: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/ChangeTaskName", fmt.Sprintf(`{"task_id":"%s", "name":"%s"}`, taskID, "Update Task Name"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// CompleteTask: TaskIDが空の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// CompleteTask: TaskIDが長すぎる場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, strings.Repeat("*", 51)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// CompleteTask: 存在しないTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, "another"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 404, res.status, "404エラーになること")

	// CompleteTask: 他人のTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, anotherTaskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 403, res.status, "パーミッションエラーになること")

	// CompleteTask: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/CompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, taskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// UncompleteTask: TaskIDが空の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/UncompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// UncompleteTask: TaskIDが長すぎる場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/UncompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, strings.Repeat("*", 51)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// UncompleteTask: 存在しないTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/UncompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, "another"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 404, res.status, "404エラーになること")

	// UncompleteTask: 他人のTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/UncompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, anotherTaskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 403, res.status, "パーミッションエラーになること")

	// UncompleteTask: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/UncompleteTask", fmt.Sprintf(`{"task_id":"%s"}`, taskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")

	// DeleteTask: TaskIDが空の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/DeleteTask", fmt.Sprintf(`{"task_id":"%s"}`, ""))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// DeleteTask: TaskIDが長すぎる場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/DeleteTask", fmt.Sprintf(`{"task_id":"%s"}`, strings.Repeat("*", 51)))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 400, res.status, "入力エラーになること")

	// DeleteTask: 存在しないTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/DeleteTask", fmt.Sprintf(`{"task_id":"%s"}`, "another"))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 404, res.status, "404エラーになること")

	// DeleteTask: 他人のTaskIDの場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/DeleteTask", fmt.Sprintf(`{"task_id":"%s"}`, anotherTaskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 403, res.status, "パーミッションエラーになること")

	// DeleteTask: 正しい入力の場合
	res, err = ts.sendPostRequest(t, token, "/rpc.task.v1.TaskService/DeleteTask", fmt.Sprintf(`{"task_id":"%s"}`, taskID))
	require.NoError(t, err, "エラーが発生しないこと")
	require.Equal(t, 200, res.status, "ステータスコードが正常であること")
}
