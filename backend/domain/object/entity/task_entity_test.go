package entity

import (
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/domain"
	"github.com/7oh2020/connect-tasklist/backend/domain/object/value"
	"github.com/stretchr/testify/require"
)

func TestTaskEntity_Validate(tt *testing.T) {
	testcases := []struct {
		title string
		arg   *Task
		err   error
	}{
		{"正常系: 正しい入力の場合", &Task{ID: value.NewID("id"), UserID: value.NewID("uid"), Name: "task"}, nil},
		{"準正常系: IDが空の場合", &Task{ID: value.NewID(""), UserID: value.NewID("uid"), Name: "task"}, &domain.ErrValidationFailed{Msg: "id is empty"}},
		{"準正常系: UserIDが空の場合", &Task{ID: value.NewID("id"), UserID: value.NewID(""), Name: "task"}, &domain.ErrValidationFailed{Msg: "id is empty"}},
		{"準正常系: nameが空の場合", &Task{ID: value.NewID("id"), UserID: value.NewID("uid"), Name: ""}, &domain.ErrValidationFailed{Msg: "name is empty"}},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			err := v.arg.Validate()

			if v.err == nil {
				require.NoError(t, err, "エラーが発生しないこと")
			} else {
				require.EqualError(t, err, v.err.Error(), "エラーが一致すること")
			}
		})
	}
}
