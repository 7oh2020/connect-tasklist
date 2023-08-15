package dto

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTaskParams_Validate(tt *testing.T) {
	testcases := []struct {
		title string
		arg   *CreateTaskParams
		err   error
	}{
		{"正常系: 正しい入力の場合", NewCreateTaskParams("uid", "task"), nil},
		{"準正常系: UserIDが半角50文字を超える場合", NewCreateTaskParams(strings.Repeat("*", 51), "title"), errors.New("id must be 50 characters or less")},
		{"準正常系: UserIDが全角50文字を超える場合", NewCreateTaskParams(strings.Repeat("あ", 51), "title"), errors.New("id must be 50 characters or less")},
		{"準正常系: Nameが半角100文字を超える場合", NewCreateTaskParams("uid", strings.Repeat("*", 101)), errors.New("name must be 100 characters or less")},
		{"準正常系: Nameが全角100文字を超える場合", NewCreateTaskParams("uid", strings.Repeat("あ", 101)), errors.New("name must be 100 characters or less")},
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
