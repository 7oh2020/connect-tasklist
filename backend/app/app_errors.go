package app

// 入力データのバリデーションに失敗した場合のエラー
type ErrInputValidationFailed struct {
	Msg string
}

func (e *ErrInputValidationFailed) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "failed to validate input"
}

// ログインに失敗した場合のエラー
type ErrLoginFailed struct {
	Msg string
}

func (e *ErrLoginFailed) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "failed to login"
}

// アプリ内部のエラー
type ErrInternal struct {
	Msg string
}

func (e *ErrInternal) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "internal error"
}
