package contextkey

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextReader_GetUserID(tt *testing.T) {
	uid := "uid"
	ctx := context.WithValue(context.Background(), ContextKeyUserID, uid)

	testcases := []struct {
		title string
		ctx   context.Context
		res   string
		err   error
	}{
		{title: "正常系: コンテキストがセットされている場合", ctx: ctx, res: uid, err: nil},
		{title: "準正常系: コンテキストがセットされていない場合", ctx: context.Background(), res: "", err: errors.New("error: context value not found for user-id")},
	}
	for _, tc := range testcases {
		tt.Run(tc.title, func(t *testing.T) {
			cr := NewContextReader()
			res, err := cr.GetUserID(tc.ctx)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.res, res)
		})

	}
}
