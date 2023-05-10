package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenManager_NewTokenManager(tt *testing.T) {
	issuer := "issuer"
	keyPath := "./test/id_rsa"

	testcases := []struct {
		title   string
		issuer  string
		keyPath string
		err     error
	}{
		{title: "正常系: 正しい入力の場合", issuer: issuer, keyPath: keyPath, err: nil},
		{title: "準正常系: issuerが空の場合", issuer: "", keyPath: keyPath, err: errors.New("error: invalid parameter")},
		{title: "準正常系: keyPathが空の場合", issuer: issuer, keyPath: "", err: errors.New("error: invalid parameter")},
		{title: "準正常系: ファイルが存在しない場合", issuer: issuer, keyPath: "undefined", err: errors.New("error: failed to read private-key-file")},
		{title: "準正常系: ファイルがPEM形式ではない場合", issuer: issuer, keyPath: "./test/id_rsa.dummy", err: errors.New("error: failed to parse pem file")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			res, err := NewTokenManager(tc.issuer, tc.keyPath)
			require.Equal(t, tc.err, err)

			if err == nil {
				require.Equal(t, issuer, res.issuer)
			}
		})
	}
}

func TestTokenManager_CreateToken(tt *testing.T) {
	keyPath := "./test/id_rsa"
	issuer := "issuer"
	uid := "uid"
	duration := 1 * time.Hour
	tm, err0 := NewTokenManager(issuer, keyPath)
	require.NoError(tt, err0)

	testcases := []struct {
		title    string
		userID   string
		duration time.Duration
		err1     error
		res      string
		err2     error
	}{
		{title: "正常系: 正しい入力の場合", userID: uid, duration: duration, err1: nil, res: uid, err2: nil},
		{title: "準正常系: userIDが空の場合", userID: "", duration: duration, err1: errors.New("error: invalid token parameter"), res: "", err2: nil},
		{title: "準正常系: トークンが有効期限切れの場合", userID: uid, duration: -1 * time.Hour, err1: nil, res: "", err2: errors.New("error: failed to verify token")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			// トークンを生成する
			token, err1 := tm.CreateToken(tc.userID, tc.duration)
			t.Logf("created token: %s", token)
			require.Equal(t, tc.err1, err1)

			// トークンからuserIDを取得する
			if err1 == nil {
				res, err2 := tm.GetUserID(token)
				require.Equal(t, tc.err2, err2)
				require.Equal(t, tc.res, res)
			}
		})
	}

}
