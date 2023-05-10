package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/require"
)

func TestTaskService_NewTaskService(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ ITaskService = (*TaskService)(nil)
	})
}

func TestTaskService_FindTaskByID(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := "id"
	task := &entity.Task{
		ID:          id,
		UserID:      "uid",
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTaskByID", ctx, id).Return(task, nil)

	testcases := []struct {
		title string
		id    string
		res   *entity.Task
		err   error
	}{
		{title: "正常系: 正しい入力の場合", id: id, res: task, err: nil},
		{title: "準正常系: idが空の場合", id: "", res: nil, err: errors.New("error: invalid parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			res, err := srv.FindTaskByID(ctx, tc.id)
			require.Equal(t, tc.err, err)

			if res != nil {
				repo.AssertExpectations(t)
				require.Equal(t, tc.res.ID, res.ID)
				require.Equal(t, tc.res.UserID, res.UserID)
				require.Equal(t, tc.res.Name, res.Name)
				require.Equal(t, tc.res.IsCompleted, res.IsCompleted)
				require.Equal(t, tc.res.CreatedAt, res.CreatedAt)
				require.Equal(t, tc.res.UpdatedAt, res.UpdatedAt)
			}
		})
	}

}

func TestTaskService_FindTasksByUserID(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	uid := "uid"
	tasks := []*entity.Task{
		{ID: "id1", UserID: uid, Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: "id2", UserID: uid, Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTasksByUserID", ctx, uid).Return(tasks, nil)

	testcases := []struct {
		title  string
		userID string
		res    []*entity.Task
		err    error
	}{
		{title: "正常系: 正しい入力の場合", userID: uid, res: tasks, err: nil},
		{title: "準正常系: userIDが空の場合", userID: "", res: tasks, err: errors.New("error: invalid parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			res, err := srv.FindTasksByUserID(ctx, tc.userID)
			require.Equal(t, tc.err, err)

			if res != nil {
				repo.AssertExpectations(t)
				require.ElementsMatch(t, tc.res, res)
			}
		})
	}
}

func TestTaskService_CreateTask(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := "id"
	uid := "uid"
	name := "task"
	task := &entity.Task{
		ID:          id,
		UserID:      uid,
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("CreateTask", ctx, task).Return(id, nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		name   string
		res    string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, name: name, res: id, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, name: name, res: "", err: errors.New("error: validation failed")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", name: name, res: "", err: errors.New("error: validation failed")},
		{title: "準正常系: nameが空の場合", id: id, userID: uid, name: "", res: "", err: errors.New("error: validation failed")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			res, err := srv.CreateTask(ctx, tc.id, tc.userID, tc.name, now)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)

			if res != "" {
				repo.AssertExpectations(t)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestTaskService_ChangeTaskName(tt *testing.T) {
	now := time.Now().UTC()
	upd := now.Add(1 * time.Second)
	id := "id"
	uid := "uid"
	name := "changed"
	ctx := context.Background()
	oldTask := &entity.Task{
		ID:          id,
		UserID:      uid,
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	newTask := &entity.Task{
		ID:          oldTask.ID,
		UserID:      oldTask.UserID,
		Name:        name,
		IsCompleted: oldTask.IsCompleted,
		CreatedAt:   oldTask.CreatedAt,
		UpdatedAt:   upd,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTaskByID", ctx, id).Return(oldTask, nil)
	repo.On("UpdateTask", ctx, newTask).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		name   string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, name: name, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, name: name, err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", name: name, err: errors.New("error: invalid parameter")},
		{title: "準正常系: nameが空の場合", id: id, userID: uid, name: "", err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが一致しない場合", id: id, userID: "another", name: name, err: errors.New("error: permission denied to update task")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			err := srv.ChangeTaskName(ctx, tc.id, tc.userID, tc.name, upd)
			require.Equal(t, tc.err, err)

			if err == nil {
				repo.AssertExpectations(t)
			}
		})
	}
}

func TestTaskService_CompleteTask(tt *testing.T) {
	now := time.Now().UTC()
	upd := now.Add(1 * time.Second)
	id := "id"
	uid := "uid"
	ctx := context.Background()
	oldTask := &entity.Task{
		ID:          id,
		UserID:      uid,
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	newTask := &entity.Task{
		ID:          oldTask.ID,
		UserID:      oldTask.UserID,
		Name:        oldTask.Name,
		IsCompleted: true,
		CreatedAt:   oldTask.CreatedAt,
		UpdatedAt:   upd,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTaskByID", ctx, id).Return(oldTask, nil)
	repo.On("UpdateTask", ctx, newTask).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが一致しない場合", id: id, userID: "another", err: errors.New("error: permission denied to update task")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			err := srv.CompleteTask(ctx, tc.id, tc.userID, upd)
			require.Equal(t, tc.err, err)

			if err == nil {
				repo.AssertExpectations(t)
			}
		})
	}
}

func TestTaskService_UncompleteTask(tt *testing.T) {
	now := time.Now().UTC()
	upd := now.Add(1 * time.Second)
	id := "id"
	uid := "uid"
	ctx := context.Background()
	oldTask := &entity.Task{
		ID:          id,
		UserID:      uid,
		Name:        "task",
		IsCompleted: true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	newTask := &entity.Task{
		ID:          oldTask.ID,
		UserID:      oldTask.UserID,
		Name:        oldTask.Name,
		IsCompleted: false,
		CreatedAt:   oldTask.CreatedAt,
		UpdatedAt:   upd,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTaskByID", ctx, id).Return(oldTask, nil)
	repo.On("UpdateTask", ctx, newTask).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが一致しない場合", id: id, userID: "another", err: errors.New("error: permission denied to update task")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			err := srv.UncompleteTask(ctx, tc.id, tc.userID, upd)
			require.Equal(t, tc.err, err)

			if err == nil {
				repo.AssertExpectations(t)
			}
		})
	}
}

func TestTaskService_DeleteTask(tt *testing.T) {
	now := time.Now().UTC()
	id := "id"
	uid := "uid"
	ctx := context.Background()
	task := &entity.Task{
		ID:          id,
		UserID:      uid,
		Name:        "task",
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo := new(mocks.ITaskRepository)
	repo.On("FindTaskByID", ctx, id).Return(task, nil)
	repo.On("DeleteTask", ctx, id).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid parameter")},
		{title: "準正常系: userIDが一致しない場合", id: id, userID: "another", err: errors.New("error: permission denied to delete task")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			srv := NewTaskService(repo)
			err := srv.DeleteTask(ctx, tc.id, tc.userID)
			require.Equal(t, tc.err, err)

			if err == nil {
				repo.AssertExpectations(t)
			}
		})
	}
}
