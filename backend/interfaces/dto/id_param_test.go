package dto

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIDParam_Validate(tt *testing.T) {
	testcases := []struct {
		title string
		arg   *IDParam
		err   error
	}{
		{"正常系: 正しい入力の場合", NewIDParam("id"), nil},
		{"準正常系: IDが半角50文字を超える場合", NewIDParam(strings.Repeat("*", 51)), errors.New("id must be 50 characters or less")},
		{"準正常系: IDが全角50文字を超える場合", NewIDParam(strings.Repeat("あ", 51)), errors.New("id must be 50 characters or less")},
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
