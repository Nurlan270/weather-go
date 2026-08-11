package request

import (
	"context"

	"github.com/Nurlan270/weather-go/internal/entity"
)

type contextKey struct {
	name string
}

var UserCtxKey = &contextKey{"user"}

func GetUserFromCtx(ctx context.Context) *entity.User {
	l, ok := ctx.Value(UserCtxKey).(*entity.User)
	if !ok {
		l = nil
	}

	return l
}
