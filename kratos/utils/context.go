package utils

import (
	"context"
	"explorer/internal/models/users"
)

type ContextKey string

func SetUserAuthenticationToContext(ctx context.Context, ua *users.UserAuthentication) context.Context {
	return context.WithValue(ctx, ContextKey("ctx_user"), ua)
}

func GetUserAuthenticationFromContext(ctx context.Context) *users.UserAuthentication {
	ua, ok := ctx.Value(ContextKey("ctx_user")).(*users.UserAuthentication)
	if !ok {
		return nil
	}
	return ua
}
