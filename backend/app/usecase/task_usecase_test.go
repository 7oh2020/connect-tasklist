package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestTaskUsecase_NewTaskUsecase(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IUserUsecase = (*UserUsecase)(nil)
	})
}

func TestTaskUsecase_FindTasksByUserID(tt *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	uid := "uid"
	tasks := []*entity.Task{
		{ID: value.NewID("t1"), UserID: value.NewID(uid), Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: value.NewID("t2"), UserID: value.NewID(uid), Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		srv := new(mocks.ITaskService)
		srv.On("FindTasksByUserID", ctx, uid).Return(tasks, nil)
		uc := NewTaskUsecase(srv)
		ret, err := uc.FindTasksByUserID(ctx, dto.NewIDParam(uid))

		require.NoError(t, err, "エラーが発生しないこと")
		require.ElementsMatch(t, tasks, ret)
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
		uid := strings.Repeat("*", 51)
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		_, err := uc.FindTasksByUserID(ctx, dto.NewIDParam(uid))

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)
	})
}

func TestTaskUsecase_CreateTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := dto.NewCreateTaskParams(uid, "task")
		srv := new(mocks.ITaskService)
		srv.On("CreateTask", ctx, uid, arg.Name()).Return(id, nil)
		uc := NewTaskUsecase(srv)
		ret, err := uc.CreateTask(ctx, arg)

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret)
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "name must be 100 characters or less"}
		arg := dto.NewCreateTaskParams(uid, strings.Repeat("*", 101))
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		_, err := uc.CreateTask(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)

	})
}

func TestTaskUsecase_ChangeTaskName(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := dto.NewChangeTaskNameParams(id, uid, "new task")
		srv := new(mocks.ITaskService)
		srv.On("ChangeTaskName", ctx, id, uid, arg.Name()).Return(nil)
		uc := NewTaskUsecase(srv)
		err := uc.ChangeTaskName(ctx, arg)

		require.NoError(t, err, "エラーが発生しないこと")
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "name must be 100 characters or less"}
		arg := dto.NewChangeTaskNameParams(id, uid, strings.Repeat("*", 101))
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		err := uc.ChangeTaskName(ctx, arg)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)

	})
}

func TestTaskUsecase_DeleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		srv := new(mocks.ITaskService)
		srv.On("DeleteTask", ctx, id, uid).Return(nil)
		uc := NewTaskUsecase(srv)
		err := uc.DeleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.NoError(t, err, "エラーが発生しないこと")
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
		id := strings.Repeat("*", 101)
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		err := uc.DeleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)
	})
}

func TestTaskUsecase_CompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		srv := new(mocks.ITaskService)
		srv.On("CompleteTask", ctx, id, uid).Return(nil)
		uc := NewTaskUsecase(srv)
		err := uc.CompleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.NoError(t, err, "エラーが発生しないこと")
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
		id := strings.Repeat("*", 101)
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		err := uc.CompleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)
	})
}

func TestTaskUsecase_UncompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		srv := new(mocks.ITaskService)
		srv.On("UncompleteTask", ctx, id, uid).Return(nil)
		uc := NewTaskUsecase(srv)
		err := uc.UncompleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.NoError(t, err, "エラーが発生しないこと")
		srv.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &app.ErrInputValidationFailed{Msg: "id must be 50 characters or less"}
		id := strings.Repeat("*", 101)
		srv := new(mocks.ITaskService)
		uc := NewTaskUsecase(srv)
		err := uc.UncompleteTask(ctx, dto.NewIDParam(id), dto.NewIDParam(uid))

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		srv.AssertExpectations(t)
	})
}
