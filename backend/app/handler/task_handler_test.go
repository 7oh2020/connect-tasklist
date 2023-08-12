package handler

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	taskv1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1/taskv1connect"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTaskHandler_NewTaskHandler(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ taskv1connect.TaskServiceHandler = (*TaskHandler)(nil)
	})
}

func TestTaskHandler_GetTaskList(tt *testing.T) {
	now := time.Now().UTC()
	uid := "uid"
	ctx := context.Background()
	ent := []*entity.Task{
		{ID: "id1", UserID: uid, Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: "id2", UserID: uid, Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}
	tasks := []*taskv1.Task{
		{Id: ent[0].ID, UserId: ent[0].UserID, Name: ent[0].Name, IsCompleted: ent[0].IsCompleted, CreatedAt: timestamppb.New(ent[0].CreatedAt), UpdatedAt: timestamppb.New(ent[0].UpdatedAt)},
		{Id: ent[1].ID, UserId: ent[1].UserID, Name: ent[1].Name, IsCompleted: ent[1].IsCompleted, CreatedAt: timestamppb.New(ent[1].CreatedAt), UpdatedAt: timestamppb.New(ent[1].UpdatedAt)},
	}

	im := new(mocks.IIDManager)
	cm := new(mocks.IClockManager)
	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)

	uct := new(mocks.ITaskUsecase)
	uct.On("GetTaskList", ctx, uid).Return(ent, nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.GetTaskListRequest]
		res   *connect.Response[taskv1.GetTaskListResponse]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.GetTaskListRequest{}), res: connect.NewResponse(&taskv1.GetTaskListResponse{Tasks: tasks}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			res, err := hdr.GetTaskList(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestTaskHandler_CreateTask(tt *testing.T) {
	id := "id"
	uid := "uid"
	name := "task"
	now := time.Now().UTC()
	ctx := context.Background()

	im := new(mocks.IIDManager)
	im.On("GenerateID").Return(id)
	cm := new(mocks.IClockManager)
	cm.On("GetNow").Return(now)
	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)

	uct := new(mocks.ITaskUsecase)
	uct.On("CreateTask", ctx, id, uid, name, now).Return(id, nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.CreateTaskRequest]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.CreateTaskRequest{Name: name}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			_, err := hdr.CreateTask(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
			}
		})
	}
}

func TestTaskHandler_ChangeTaskName(tt *testing.T) {
	id := "id"
	uid := "uid"
	name := "task"
	now := time.Now().UTC()
	ctx := context.Background()

	im := new(mocks.IIDManager)
	cm := new(mocks.IClockManager)
	cm.On("GetNow").Return(now)
	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)

	uct := new(mocks.ITaskUsecase)
	uct.On("ChangeTaskName", ctx, id, uid, name, now).Return(nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.ChangeTaskNameRequest]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.ChangeTaskNameRequest{Id: id, Name: name}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			_, err := hdr.ChangeTaskName(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
			}
		})
	}
}

func TestTaskHandler_CompleteTask(tt *testing.T) {
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	ctx := context.Background()

	im := new(mocks.IIDManager)
	cm := new(mocks.IClockManager)
	cm.On("GetNow").Return(now)
	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)

	uct := new(mocks.ITaskUsecase)
	uct.On("CompleteTask", ctx, id, uid, now).Return(nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.CompleteTaskRequest]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.CompleteTaskRequest{TaskId: id}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			_, err := hdr.CompleteTask(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
			}
		})
	}
}

func TestTaskHandler_UncompleteTask(tt *testing.T) {
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	ctx := context.Background()

	im := new(mocks.IIDManager)
	cm := new(mocks.IClockManager)
	cm.On("GetNow").Return(now)
	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)

	uct := new(mocks.ITaskUsecase)
	uct.On("UncompleteTask", ctx, id, uid, now).Return(nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.UncompleteTaskRequest]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.UncompleteTaskRequest{TaskId: id}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			_, err := hdr.UncompleteTask(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
			}
		})
	}
}

func TestTaskHandler_DeleteTask(tt *testing.T) {
	id := "id"
	uid := "uid"
	ctx := context.Background()

	im := new(mocks.IIDManager)

	cr := new(mocks.IContextReader)
	cr.On("GetUserID", ctx).Return(uid, nil)
	cm := new(mocks.IClockManager)
	uct := new(mocks.ITaskUsecase)
	uct.On("DeleteTask", ctx, id, uid).Return(nil)

	testcases := []struct {
		title string
		ctx   context.Context
		arg   *connect.Request[taskv1.DeleteTaskRequest]
		err   error
	}{
		{title: "正常系: 正しい入力の場合", ctx: ctx, arg: connect.NewRequest(&taskv1.DeleteTaskRequest{TaskId: id}), err: nil},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			hdr := NewTaskHandler(im, cm, cr, uct)
			_, err := hdr.DeleteTask(tc.ctx, tc.arg)
			require.Equal(t, tc.err, err)

			if err == nil {
				im.AssertExpectations(t)
				cr.AssertExpectations(t)
				uct.AssertExpectations(t)
			}
		})
	}
}
