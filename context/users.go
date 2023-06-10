package context

import (
	"context"
	"lenslocked/domain/entity"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *entity.User {
	if v := ctx.Value(userKey); v != nil {
		if user, ok := v.(*entity.User); ok {
			return user
		}
	}
	return nil
}
