package service

import (
	"context"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestTaskService_NewTaskService(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ ITaskService = (*TaskService)(nil)
	})
}

func TestTaskService_FindTaskByID(tt *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	id := "id"
	uid := "uid"
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		ret, err := srv.FindTaskByID(ctx, id)

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret.ID.Value())
		require.Equal(t, task.UserID.Value(), ret.UserID.Value())
		require.Equal(t, task.Name, ret.Name)
		require.Equal(t, task.IsCompleted, ret.IsCompleted)
		require.Equal(t, task.CreatedAt, ret.CreatedAt)
		require.Equal(t, task.UpdatedAt, ret.UpdatedAt)
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: TaskIDが空の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		id := ""
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.FindTaskByID(ctx, id)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないTaskIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "task not found"}
		id := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.FindTaskByID(ctx, id)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_FindTasksByUserID(tt *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	uid := "uid"
	tasks := []*entity.Task{
		{ID: value.NewID("t1"), UserID: value.NewID(uid), Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: value.NewID("t2"), UserID: value.NewID(uid), Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		repo := new(mocks.ITaskRepository)
		repo.On("FindTasksByUserID", ctx, uid).Return(tasks, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		ret, err := srv.FindTasksByUserID(ctx, uid)

		require.NoError(t, err, "エラーが発生しないこと")
		require.ElementsMatch(t, tasks, ret)
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: UserIDが空の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		uid := ""
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.FindTasksByUserID(ctx, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないUserIDの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		uid := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTasksByUserID", ctx, uid).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.FindTasksByUserID(ctx, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_CreateTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		repo := new(mocks.ITaskRepository)
		repo.On("CreateTask", ctx, task).Return(id, nil)
		im := new(mocks.IIDManager)
		im.On("GenerateID").Return(id)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(now)
		srv := NewTaskService(repo, im, cm)
		ret, err := srv.CreateTask(ctx, uid, task.Name)

		require.NoError(t, err, "エラーが発生しないこと")
		require.Equal(t, id, ret)
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "name is empty"}
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		im.On("GenerateID").Return(id)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(now)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.CreateTask(ctx, uid, "")

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: クエリエラーの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		repo := new(mocks.ITaskRepository)
		repo.On("CreateTask", ctx, task).Return("", errExp)
		im := new(mocks.IIDManager)
		im.On("GenerateID").Return(id)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(now)
		srv := NewTaskService(repo, im, cm)
		_, err := srv.CreateTask(ctx, uid, task.Name)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_ChangeTaskName(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	upd := now.Add(time.Second)
	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        "new task",
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.ChangeTaskName(ctx, arg.ID.Value(), arg.UserID.Value(), arg.Name)

		require.NoError(t, err, "エラーが発生しないこと")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "name is empty"}
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        "",
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.ChangeTaskName(ctx, arg.ID.Value(), arg.UserID.Value(), arg.Name)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないTaskIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "task not found"}
		arg := &entity.Task{
			ID:          value.NewID("another"),
			UserID:      task.UserID,
			Name:        "new task",
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, arg.ID.Value()).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.ChangeTaskName(ctx, arg.ID.Value(), arg.UserID.Value(), arg.Name)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: アクセス権がない場合", func(t *testing.T) {
		errExp := &domain.ErrPermissionDenied{}
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      value.NewID("another"),
			Name:        "new task",
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.ChangeTaskName(ctx, arg.ID.Value(), arg.UserID.Value(), arg.Name)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: クエリエラーの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        "new task",
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.ChangeTaskName(ctx, arg.ID.Value(), arg.UserID.Value(), arg.Name)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_DeleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("DeleteTask", ctx, id).Return(nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.DeleteTask(ctx, id, uid)

		require.NoError(t, err, "エラーが発生しないこと")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		id := ""
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.DeleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないTaskIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "task not found"}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.DeleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: アクセス権がない場合", func(t *testing.T) {
		errExp := &domain.ErrPermissionDenied{}
		uid := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.DeleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: クエリエラーの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("DeleteTask", ctx, id).Return(errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.DeleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_CompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	upd := now.Add(time.Second)

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        task.Name,
			IsCompleted: true,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.CompleteTask(ctx, id, uid)

		require.NoError(t, err, "エラーが発生しないこと")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		id := ""
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.CompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないTaskIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "task not found"}
		id := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.CompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: アクセス権がない場合", func(t *testing.T) {
		errExp := &domain.ErrPermissionDenied{}
		uid := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.CompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: クエリエラーの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        task.Name,
			IsCompleted: true,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.CompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}

func TestTaskService_UncompleteTask(tt *testing.T) {
	ctx := context.Background()
	id := "id"
	uid := "uid"
	now := time.Now().UTC()
	task := &entity.Task{
		ID:          value.NewID(id),
		UserID:      value.NewID(uid),
		Name:        "task",
		IsCompleted: true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	upd := now.Add(time.Second)

	tt.Run("正常系: 正しい入力の場合", func(t *testing.T) {
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        task.Name,
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.UncompleteTask(ctx, id, uid)

		require.NoError(t, err, "エラーが発生しないこと")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 不正な入力の場合", func(t *testing.T) {
		errExp := &domain.ErrValidationFailed{Msg: "id is empty"}
		id := ""
		repo := new(mocks.ITaskRepository)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.UncompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: 存在しないTaskIDの場合", func(t *testing.T) {
		errExp := &domain.ErrNotFound{Msg: "task not found"}
		id := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(nil, errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.UncompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: アクセス権がない場合", func(t *testing.T) {
		errExp := &domain.ErrPermissionDenied{}
		uid := "another"
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		srv := NewTaskService(repo, im, cm)
		err := srv.UncompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
	tt.Run("準正常系: クエリエラーの場合", func(t *testing.T) {
		errExp := &domain.ErrQueryFailed{}
		arg := &entity.Task{
			ID:          task.ID,
			UserID:      task.UserID,
			Name:        task.Name,
			IsCompleted: false,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   upd,
		}
		repo := new(mocks.ITaskRepository)
		repo.On("FindTaskByID", ctx, id).Return(task, nil)
		repo.On("UpdateTask", ctx, arg).Return(errExp)
		im := new(mocks.IIDManager)
		cm := new(mocks.IClockManager)
		cm.On("GetNow").Return(upd)
		srv := NewTaskService(repo, im, cm)
		err := srv.UncompleteTask(ctx, id, uid)

		require.EqualError(t, err, errExp.Error(), "エラーが一致すること")
		repo.AssertExpectations(t)
		im.AssertExpectations(t)
		cm.AssertExpectations(t)
	})
}
