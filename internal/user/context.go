package user

import (
	"context"
	"net/http"

	"github.com/Nurlan270/weather-go/internal/entity"
)

type contextKey struct {
	name string
}

var userKey = &contextKey{"user"}

func NewContext(ctx context.Context, u *entity.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func FromContext(ctx context.Context) *entity.User {
	u, ok := ctx.Value(userKey).(*entity.User)
	if !ok {
		return nil
	}

	return u
}

func FromRequest(r *http.Request) *entity.User {
	return FromContext(r.Context())
}
