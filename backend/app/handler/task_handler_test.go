package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	task_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1/task_v1connect"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler_NewTaskHandler(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ task_v1connect.TaskServiceHandler = (*TaskHandler)(nil)
	})
}

func TestTaskHandler_GetTaskList(tt *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	uid := "uid"
	tasks := []*entity.Task{
		{ID: value.NewID("t1"), UserID: value.NewID(uid), Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: value.NewID("t2"), UserID: value.NewID(uid), Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}
	arg := &task_v1.GetTaskListRequest{}
	param := dto.NewIDParam(uid)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("FindTasksByUserID", ctx, param).Return(tasks, nil)
			} else {
				uc.On("FindTasksByUserID", ctx, param).Return(nil, v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			ret, err := hdr.GetTaskList(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
				for i, v := range ret.Msg.Tasks {
					require.Equal(t, tasks[i].ID.Value(), v.Id)
					require.Equal(t, tasks[i].UserID.Value(), v.UserId)
					require.Equal(t, tasks[i].Name, v.Name)
					require.Equal(t, tasks[i].IsCompleted, v.IsCompleted)
				}
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_CreateTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	arg := &task_v1.CreateTaskRequest{Name: "task"}
	param := dto.NewCreateTaskParams(uid, arg.Name)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("CreateTask", ctx, param).Return(id, nil)
			} else {
				uc.On("CreateTask", ctx, param).Return("", v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			ret, err := hdr.CreateTask(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
				require.Equal(t, id, ret.Msg.CreatedId)
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_ChangeTaskName(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	arg := &task_v1.ChangeTaskNameRequest{TaskId: id, Name: "new task"}
	param := dto.NewChangeTaskNameParams(arg.TaskId, uid, arg.Name)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: タスクが存在しない場合", &domain.ErrNotFound{}, "not_found"},
		{"準正常系: アクセス権がない場合", &domain.ErrPermissionDenied{}, "permission_denied"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("ChangeTaskName", ctx, param).Return(nil)
			} else {
				uc.On("ChangeTaskName", ctx, param).Return(v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			_, err := hdr.ChangeTaskName(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_DeleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	arg := &task_v1.DeleteTaskRequest{TaskId: id}
	paramID := dto.NewIDParam(arg.TaskId)
	paramUserID := dto.NewIDParam(uid)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: タスクが存在しない場合", &domain.ErrNotFound{}, "not_found"},
		{"準正常系: アクセス権がない場合", &domain.ErrPermissionDenied{}, "permission_denied"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("DeleteTask", ctx, paramID, paramUserID).Return(nil)
			} else {
				uc.On("DeleteTask", ctx, paramID, paramUserID).Return(v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			_, err := hdr.DeleteTask(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_CompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	arg := &task_v1.CompleteTaskRequest{TaskId: id}
	paramID := dto.NewIDParam(arg.TaskId)
	paramUserID := dto.NewIDParam(uid)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: タスクが存在しない場合", &domain.ErrNotFound{}, "not_found"},
		{"準正常系: アクセス権がない場合", &domain.ErrPermissionDenied{}, "permission_denied"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("CompleteTask", ctx, paramID, paramUserID).Return(nil)
			} else {
				uc.On("CompleteTask", ctx, paramID, paramUserID).Return(v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			_, err := hdr.CompleteTask(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestTaskHandler_UncompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	arg := &task_v1.UncompleteTaskRequest{TaskId: id}
	paramID := dto.NewIDParam(arg.TaskId)
	paramUserID := dto.NewIDParam(uid)
	req := connect.NewRequest(arg)

	testcases := []struct {
		title   string
		err     error
		codeStr string
	}{
		{"正常系: 正しい入力の場合", nil, ""},
		{"準正常系: アプリ側バリデーションエラーの場合", &app.ErrInputValidationFailed{}, "invalid_argument"},
		{"準正常系: ドメイン側バリデーションエラーの場合", &domain.ErrValidationFailed{}, "invalid_argument"},
		{"準正常系: タスクが存在しない場合", &domain.ErrNotFound{}, "not_found"},
		{"準正常系: アクセス権がない場合", &domain.ErrPermissionDenied{}, "permission_denied"},
		{"準正常系: クエリエラーの場合", &domain.ErrQueryFailed{}, "aborted"},
		{"準正常系: その他のエラーの場合", &app.ErrInternal{}, "unknown"},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			uc := new(mocks.ITaskUsecase)
			if v.err == nil {
				uc.On("UncompleteTask", ctx, paramID, paramUserID).Return(nil)
			} else {
				uc.On("UncompleteTask", ctx, paramID, paramUserID).Return(v.err)
			}
			cr := new(mocks.IContextReader)
			cr.On("GetUserID", ctx).Return(uid, nil)
			hdr := NewTaskHandler(uc, cr)
			_, err := hdr.UncompleteTask(ctx, req)

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
			} else {
				errMsg := fmt.Sprintf("%s: %s", v.codeStr, v.err.Error())
				require.EqualError(t, err, errMsg, "エラーが一致すること")
			}
			uc.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}
