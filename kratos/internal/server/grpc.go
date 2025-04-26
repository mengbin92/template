package server

import (
	v1 "explorer/api/explorer/v1"
	"explorer/internal/conf"
	"explorer/internal/middleware"
	"explorer/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, basicService *service.BasicService, userService *service.UserService, logger log.Logger) *grpc.Server {
	authMiddleware := middleware.NewAuthMiddleware(userService.UserManager, logger)
	
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			middleware.AuthMiddlewareWrap(authMiddleware),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterBasicServer(srv, basicService)
	v1.RegisterUserServer(srv, userService)
	return srv
}
