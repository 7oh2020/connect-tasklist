package domain

// パラメータのバリデーションに失敗した時のエラー
type ErrValidationFailed struct {
	Msg string
}

func (e *ErrValidationFailed) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "failed to validate parameter"
}

// 該当データが存在しない場合のエラー
type ErrNotFound struct {
	Msg string
}

func (e *ErrNotFound) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "not found"
}

// DBのクエリに失敗した場合のエラー
type ErrQueryFailed struct {
	Msg string
}

func (e *ErrQueryFailed) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "failed to query"
}

// リソースへのアクセス権がなかった場合のエラー
type ErrPermissionDenied struct {
	Msg string
}

func (e *ErrPermissionDenied) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return "permission denied"
}
