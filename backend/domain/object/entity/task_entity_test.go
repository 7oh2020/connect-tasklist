package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTaskEntity_NewTaskEntity(tt *testing.T) {
	now := time.Now().UTC()
	id := "id"
	uid := "uid"
	name := "task"
	task := &Task{
		ID:          id,
		UserID:      uid,
		Name:        name,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	testcases := []struct {
		title  string
		id     string
		userID string
		name   string
		res    *Task
		err    error
	}{
		{title: "正常系: 正しい入力の場合", id: id, userID: uid, name: name, res: task, err: nil},
		{title: "準正常系: idが空の場合", id: "", userID: uid, name: name, res: nil, err: errors.New("error: validation failed")},
		{title: "準正常系: userIDが空の場合", id: id, userID: "", name: name, res: nil, err: errors.New("error: validation failed")},
		{title: "準正常系: nameが空の場合", id: id, userID: uid, name: "", res: nil, err: errors.New("error: validation failed")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			res, err := NewTask(tc.id, tc.userID, tc.name, now)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)
		})
	}
}
