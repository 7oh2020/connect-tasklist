package contextkey

// コンテキスト値に対応するキー文字列
type contextKey string

const (
	// ユーザーを識別するID。認証後にコンテキストにセットされる
	ContextKeyUserID contextKey = "ctx-user-id"
)
