package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/interfaces/dto"
	task_v1 "github.com/7oh2020/connect-tasklist/backend/interfaces/rpc/task/v1"
	"github.com/7oh2020/connect-tasklist/backend/util/contextkey"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TaskServiceHandlerの実装
type TaskHandler struct {
	usecase.ITaskUsecase
	contextkey.IContextReader
}

func NewTaskHandler(uc usecase.ITaskUsecase, cr contextkey.IContextReader) *TaskHandler {
	return &TaskHandler{uc, cr}
}

func (h *TaskHandler) GetTaskList(ctx context.Context, arg *connect.Request[task_v1.GetTaskListRequest]) (*connect.Response[task_v1.GetTaskListResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	res, err := h.ITaskUsecase.FindTasksByUserID(ctx, dto.NewIDParam(uid))
	if err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	tasks := make([]*task_v1.Task, len(res))
	for i, v := range res {
		tasks[i] = &task_v1.Task{
			Id:          v.ID.Value(),
			UserId:      v.UserID.Value(),
			Name:        v.Name,
			IsCompleted: v.IsCompleted,
			CreatedAt:   timestamppb.New(v.CreatedAt),
			UpdatedAt:   timestamppb.New(v.UpdatedAt),
		}
	}
	return connect.NewResponse(&task_v1.GetTaskListResponse{
		Tasks: tasks,
	}), nil
}

func (h *TaskHandler) CreateTask(ctx context.Context, arg *connect.Request[task_v1.CreateTaskRequest]) (*connect.Response[task_v1.CreateTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	createdID, err := h.ITaskUsecase.CreateTask(ctx, dto.NewCreateTaskParams(uid, arg.Msg.Name))
	if err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&task_v1.CreateTaskResponse{
		CreatedId: createdID,
	}), nil
}

func (h *TaskHandler) ChangeTaskName(ctx context.Context, arg *connect.Request[task_v1.ChangeTaskNameRequest]) (*connect.Response[task_v1.ChangeTaskNameResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.ChangeTaskName(ctx, dto.NewChangeTaskNameParams(arg.Msg.TaskId, uid, arg.Msg.Name)); err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeNotFound, e)
		case *domain.ErrPermissionDenied:
			return nil, connect.NewError(connect.CodePermissionDenied, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&task_v1.ChangeTaskNameResponse{}), nil
}

func (h *TaskHandler) CompleteTask(ctx context.Context, arg *connect.Request[task_v1.CompleteTaskRequest]) (*connect.Response[task_v1.CompleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.CompleteTask(ctx, dto.NewIDParam(arg.Msg.TaskId), dto.NewIDParam(uid)); err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeNotFound, e)
		case *domain.ErrPermissionDenied:
			return nil, connect.NewError(connect.CodePermissionDenied, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&task_v1.CompleteTaskResponse{}), nil
}

func (h *TaskHandler) UncompleteTask(ctx context.Context, arg *connect.Request[task_v1.UncompleteTaskRequest]) (*connect.Response[task_v1.UncompleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.UncompleteTask(ctx, dto.NewIDParam(arg.Msg.TaskId), dto.NewIDParam(uid)); err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeNotFound, e)
		case *domain.ErrPermissionDenied:
			return nil, connect.NewError(connect.CodePermissionDenied, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&task_v1.UncompleteTaskResponse{}), nil
}

func (h *TaskHandler) DeleteTask(ctx context.Context, arg *connect.Request[task_v1.DeleteTaskRequest]) (*connect.Response[task_v1.DeleteTaskResponse], error) {
	// コンテキストから値を取得する
	var uid string
	var err error
	if uid, err = h.IContextReader.GetUserID(ctx); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if err := h.ITaskUsecase.DeleteTask(ctx, dto.NewIDParam(arg.Msg.TaskId), dto.NewIDParam(uid)); err != nil {
		switch e := err.(type) {
		case *app.ErrInputValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrValidationFailed:
			return nil, connect.NewError(connect.CodeInvalidArgument, e)
		case *domain.ErrNotFound:
			return nil, connect.NewError(connect.CodeNotFound, e)
		case *domain.ErrPermissionDenied:
			return nil, connect.NewError(connect.CodePermissionDenied, e)
		case *domain.ErrQueryFailed:
			return nil, connect.NewError(connect.CodeAborted, e)
		default:
			return nil, connect.NewError(connect.CodeUnknown, e)
		}
	}
	return connect.NewResponse(&task_v1.DeleteTaskResponse{}), nil

}
