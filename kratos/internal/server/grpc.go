package server

import (
	v1 "explorer/api/explorer/v1"
	"explorer/internal/conf"
	"explorer/internal/middleware"
	"explorer/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"google.golang.org/grpc"
	krpc "github.com/go-kratos/kratos/v2/transport/grpc"
)

var (
	// standardGrpcServer is the underlying grpc.Server instance used for gRPC-Web
	standardGrpcServer *grpc.Server
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, basicService *service.BasicService, userService *service.UserService, logger log.Logger) *krpc.Server {
	authMiddleware := middleware.NewAuthMiddleware(userService.UserManager, logger)

	var opts = []krpc.ServerOption{
		krpc.Middleware(
			recovery.Recovery(),
			middleware.AuthMiddlewareWrap(authMiddleware),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, krpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, krpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, krpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := krpc.NewServer(opts...)
	v1.RegisterBasicServer(srv, basicService)
	v1.RegisterUserServer(srv, userService)
	return srv
}

// GetStandardGRPCServer returns the underlying standard grpc.Server instance.
// This is used by the gRPC-Web wrapper to handle gRPC-Web requests.
func GetStandardGRPCServer() *grpc.Server {
	return standardGrpcServer
}