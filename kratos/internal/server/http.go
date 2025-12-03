// Package server provides server initialization for both gRPC and HTTP servers.
package server

import (
	v1 "kratos-project-template/api/demo/v1"
	"kratos-project-template/internal/conf"
	"kratos-project-template/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// NewHTTPServer creates and configures a new HTTP server instance.
func NewHTTPServer(c *conf.Server, logger log.Logger) *khttp.Server {
	// Configure CORS
	var corsOpts []handlers.CORSOption
	allowedOrigins := []string{"*"}
	allowCredentials := false

	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		allowCredentials = false
	}

	corsOpts = append(corsOpts,
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		}),
		handlers.MaxAge(86400),
	)

	if allowCredentials {
		corsOpts = append(corsOpts, handlers.AllowCredentials())
	}

	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
	}

	corsMiddleware := handlers.CORS(corsOpts...)

	opts = append(opts, khttp.Filter(corsMiddleware))

	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)

	demoService := service.NewDemoService()
	v1.RegisterDemoHTTPServer(srv, demoService)

	return srv
}
