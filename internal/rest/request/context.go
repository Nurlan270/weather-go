package request

import "context"

type contextKey struct {
	name string
}

var LoginCtxKey = &contextKey{"login"}

func GetLoginFromCtx(ctx context.Context) string {
	l, ok := ctx.Value(LoginCtxKey).(string)
	if !ok {
		l = ""
	}

	return l
}
