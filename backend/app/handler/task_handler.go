package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/app/util/clock"
	"github.com/7oh2020/connect-tasklist/backend/app/util/contextkey"
	"github.com/7oh2020/connect-tasklist/backend/app/util/identification"
	taskv1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TaskServiceHandlerの実装
type TaskHandler struct {
	identification.IIDManager
	clock.IClockManager
	contextkey.IContextReader
	usecase.ITaskUsecase
}

func NewTaskHandler(im identification.IIDManager, cm clock.IClockManager, cr contextkey.IContextReader, uct usecase.ITaskUsecase) *TaskHandler {
	return &TaskHandler{im, cm, cr, uct}
}

func (h *TaskHandler) GetTaskList(ctx context.Context, arg *connect.Request[taskv1.GetTaskListRequest]) (*connect.Response[taskv1.GetTaskListResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	res, err := h.ITaskUsecase.GetTaskList(ctx, uid)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	tasks := make([]*taskv1.Task, len(res))
	for i, v := range res {
		tasks[i] = &taskv1.Task{
			Id:          v.ID,
			UserId:      v.UserID,
			Name:        v.Name,
			IsCompleted: v.IsCompleted,
			CreatedAt:   timestamppb.New(v.CreatedAt),
			UpdatedAt:   timestamppb.New(v.UpdatedAt),
		}
	}
	return connect.NewResponse(&taskv1.GetTaskListResponse{
		Tasks: tasks,
	}), nil
}

func (h *TaskHandler) CreateTask(ctx context.Context, arg *connect.Request[taskv1.CreateTaskRequest]) (*connect.Response[taskv1.CreateTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	createdID, err := h.ITaskUsecase.CreateTask(ctx, h.IIDManager.GenerateID(), uid, arg.Msg.Name, h.IClockManager.GetNow())
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&taskv1.CreateTaskResponse{
		CreatedId: createdID,
	}), nil
}

func (h *TaskHandler) ChangeTaskName(ctx context.Context, arg *connect.Request[taskv1.ChangeTaskNameRequest]) (*connect.Response[taskv1.ChangeTaskNameResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.ChangeTaskName(ctx, arg.Msg.Id, uid, arg.Msg.Name, h.IClockManager.GetNow()); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&taskv1.ChangeTaskNameResponse{}), nil
}

func (h *TaskHandler) CompleteTask(ctx context.Context, arg *connect.Request[taskv1.CompleteTaskRequest]) (*connect.Response[taskv1.CompleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.CompleteTask(ctx, arg.Msg.TaskId, uid, h.IClockManager.GetNow()); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&taskv1.CompleteTaskResponse{}), nil
}

func (h *TaskHandler) UncompleteTask(ctx context.Context, arg *connect.Request[taskv1.UncompleteTaskRequest]) (*connect.Response[taskv1.UncompleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.UncompleteTask(ctx, arg.Msg.TaskId, uid, h.IClockManager.GetNow()); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&taskv1.UncompleteTaskResponse{}), nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, arg *connect.Request[taskv1.DeleteTaskRequest]) (*connect.Response[taskv1.DeleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.DeleteTask(ctx, arg.Msg.TaskId, uid); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}
	return connect.NewResponse(&taskv1.DeleteTaskResponse{}), nil

}
