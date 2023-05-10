package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/7oh2020/connect-tasklist/backend/domain/object/entity"
	"github.com/7oh2020/connect-tasklist/backend/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskUsecase_NewTaskUsecase(tt *testing.T) {
	tt.Run("異常系: structがinterfaceを実装しているか", func(t *testing.T) {
		var _ IUserUsecase = (*UserUsecase)(nil)
	})
}

func TestTaskUsecase_GetTaskList(tt *testing.T) {
	ctx := context.Background()
	uid := strings.Repeat("*", 50)
	now := time.Now().UTC()
	tasks := []*entity.Task{
		{ID: "id1", UserID: uid, Name: "task1", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
		{ID: "id2", UserID: uid, Name: "task2", IsCompleted: false, CreatedAt: now, UpdatedAt: now},
	}

	srv := new(mocks.ITaskService)
	srv.On("FindTasksByUserID", ctx, uid).Return(tasks, nil)

	testcases := []struct {
		title  string
		userID string
		res    []*entity.Task
		err    error
	}{
		{title: "正常系: 正しい入力の場合", userID: uid, res: tasks, err: nil},
		{title: "準正常系: userIDが空の場合", userID: "", res: nil, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", userID: strings.Repeat("a", 51), res: nil, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", userID: strings.Repeat("あ", 51), res: nil, err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			res, err := uc.GetTaskList(ctx, tc.userID)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)
			if res != nil {
				srv.AssertExpectations(t)
				assert.ElementsMatch(t, tasks, res, "要素が一致すること")
			}
		})
	}
}

func TestTaskUsecase_CreateTask(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	uid := strings.Repeat("*", 50)
	name := strings.Repeat("*", 100)

	srv := new(mocks.ITaskService)
	srv.On("CreateTask", ctx, id, uid, name, now).Return(id, nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		name   string
		res    string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, name: name, res: id, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), userID: uid, name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), userID: uid, name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", id: id, userID: strings.Repeat("a", 51), name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", id: id, userID: strings.Repeat("あ", 51), name: name, res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが空の場合", id: id, userID: uid, name: "", res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが半角101文字の場合", id: id, userID: uid, name: strings.Repeat("a", 101), res: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが全角101文字の場合", id: id, userID: uid, name: strings.Repeat("あ", 101), res: "", err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			res, err := uc.CreateTask(ctx, tc.id, tc.userID, tc.name, now)
			require.Equal(t, tc.err, err)

			if res != "" {
				srv.AssertExpectations(t)
				require.Equal(t, tc.res, res)
			}
		})
	}

}

func TestTaskUsecase_ChangeTaskName(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	uid := strings.Repeat("*", 50)
	name := strings.Repeat("*", 100)

	srv := new(mocks.ITaskService)
	srv.On("ChangeTaskName", ctx, id, uid, name, now).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		name   string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, name: name, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), userID: uid, name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), userID: uid, name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", id: id, userID: strings.Repeat("a", 51), name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", id: id, userID: strings.Repeat("あ", 51), name: name, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが空の場合", id: id, userID: uid, name: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが半角101文字の場合", id: id, userID: uid, name: strings.Repeat("a", 101), err: errors.New("error: invalid request parameter")},
		{title: "準正常系: nameが全角101文字の場合", id: id, userID: uid, name: strings.Repeat("あ", 101), err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			err := uc.ChangeTaskName(ctx, tc.id, tc.userID, tc.name, now)
			require.Equal(t, tc.err, err)

			if err == nil {
				srv.AssertExpectations(t)
			}
		})
	}
}

func TestTaskUsecase_CompleteTask(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	uid := strings.Repeat("*", 50)

	srv := new(mocks.ITaskService)
	srv.On("CompleteTask", ctx, id, uid, now).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", id: id, userID: strings.Repeat("a", 51), err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", id: id, userID: strings.Repeat("あ", 51), err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			err := uc.CompleteTask(ctx, tc.id, tc.userID, now)
			require.Equal(t, tc.err, err)

			if err == nil {
				srv.AssertExpectations(t)
			}
		})
	}
}

func TestTaskUsecase_UncompleteTask(tt *testing.T) {
	now := time.Now().UTC()
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	uid := strings.Repeat("*", 50)

	srv := new(mocks.ITaskService)
	srv.On("UncompleteTask", ctx, id, uid, now).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", id: id, userID: strings.Repeat("a", 51), err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", id: id, userID: strings.Repeat("あ", 51), err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			err := uc.UncompleteTask(ctx, tc.id, tc.userID, now)
			require.Equal(t, tc.err, err)

			if err == nil {
				srv.AssertExpectations(t)
			}
		})
	}
}

func TestTaskUsecase_DeleteTask(tt *testing.T) {
	ctx := context.Background()
	id := strings.Repeat("*", 50)
	uid := strings.Repeat("*", 50)

	srv := new(mocks.ITaskService)
	srv.On("DeleteTask", ctx, id, uid).Return(nil)

	testcases := []struct {
		title  string
		id     string
		userID string
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが半角51文字の場合", id: strings.Repeat("a", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: idが全角51文字の場合", id: strings.Repeat("あ", 51), userID: uid, err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが半角51文字の場合", id: id, userID: strings.Repeat("a", 51), err: errors.New("error: invalid request parameter")},
		{title: "準正常系: userIDが全角51文字の場合", id: id, userID: strings.Repeat("あ", 51), err: errors.New("error: invalid request parameter")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			uc := NewTaskUsecase(srv)
			err := uc.DeleteTask(ctx, tc.id, tc.userID)
			require.Equal(t, tc.err, err)

			if err == nil {
				srv.AssertExpectations(t)
			}
		})
	}
}
