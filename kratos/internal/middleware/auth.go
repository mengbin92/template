package middleware

import (
	"context"
	"explorer/internal/models/users"
	"explorer/utils"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/pkg/errors"
)

type AuthMiddleware struct {
	manager *users.UserManager
	logger *log.Helper
}

func NewAuthMiddleware(manager *users.UserManager, logger log.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		manager: manager,
		logger: log.NewHelper(log.With(logger, "module", "middleware/auth")),
	}
}

func (m *AuthMiddleware) Handler(ctx context.Context, req interface{}, next middleware.Handler) (interface{}, error) {
	tr,ok := transport.FromServerContext(ctx)
	if !ok{
		return nil, errors.New("no transport found")
	}

	m.logger.Infof("transport operation: %s", tr.Operation())

	// 跳过注册和登陆
	if tr.Operation() == "/api.explorer.v1.User/Register" || tr.Operation() == "/api.explorer.v1.User/Login" || tr.Operation() == "/api.explorer.v1.Basic/Ping" {
		return next(ctx, req)
	}

	authHeader := tr.RequestHeader().Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no auth header found")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		return nil, errors.New("no token found")
	}
	ua,err := m.manager.ValidToken(ctx, tokenStr)
	if err != nil{
		return nil, errors.Wrap(err, "get user id from token failed")
	}
	ctx = utils.SetUserAuthenticationToContext(ctx, ua)
	return next(ctx, req)
}

func AuthMiddlewareWrap(auth *AuthMiddleware) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return auth.Handler(ctx, req, handler)
		}
	}
}