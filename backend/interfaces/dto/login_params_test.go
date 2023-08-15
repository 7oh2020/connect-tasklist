package dto

import (
	"strings"
	"testing"

	"github.com/7oh2020/connect-tasklist/backend/app"
	"github.com/stretchr/testify/require"
)

func TestLoginParams_Validate(tt *testing.T) {
	testcases := []struct {
		title string
		arg   *LoginParams
		err   error
	}{
		{"正常系: 正しい入力の場合", NewLoginParams("test@example.com", "pass"), nil},
		{"準正常系: Emailが半角100文字を超える場合", NewLoginParams(strings.Repeat("*", 101), "pass"), &app.ErrInputValidationFailed{Msg: "email must be 100 characters or less"}},
		{"準正常系: Emailが全角100文字を超える場合", NewLoginParams(strings.Repeat("あ", 101), "pass"), &app.ErrInputValidationFailed{Msg: "email must be 100 characters or less"}},
		{"準正常系: Passwordが半角100文字を超える場合", NewLoginParams("test@example.com", strings.Repeat("*", 101)), &app.ErrInputValidationFailed{Msg: "password must be 100 characters or less"}},
		{"準正常系: Passwordが全角100文字を超える場合", NewLoginParams("test@example.com", strings.Repeat("あ", 101)), &app.ErrInputValidationFailed{Msg: "password must be 100 characters or less"}},
		{"準正常系: Emailのフォーマットが不正な場合", NewLoginParams("email", "pass"), &app.ErrInputValidationFailed{Msg: "invalid email"}},
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
