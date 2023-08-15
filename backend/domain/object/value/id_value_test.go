package value

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestID_Validate(tt *testing.T) {
	testcases := []struct {
		title string
		arg   *ID
		err   error
	}{
		{"正常系: 入力データが正しい場合", NewID("id"), nil},
		{"準正常系: 入力データが空の場合", NewID(""), errors.New("id is empty")},
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

func TestID_Equal(tt *testing.T) {
	testcases := []struct {
		title string
		from  string
		to    string
		ret   bool
	}{
		{"正常系: 一致する場合", "aaa", "aaa", true},
		{"正常系: 一致しない場合", "aaa", "bbb", false},
	}
	for _, v := range testcases {
		tt.Run(v.title, func(t *testing.T) {
			ret := NewID(v.from).Equal(v.to)
			require.Equal(t, v.ret, ret, "期待通りの値であること")
		})
	}
}
