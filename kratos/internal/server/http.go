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
// It sets up middleware, network configuration, address, and timeout from the provided configuration.
// The server registers the demo service HTTP handler.
//
// Parameters:
//   - c: Server configuration containing HTTP settings
//   - logger: Logger instance for server logging
//
// Returns:
//   - *khttp.Server: A configured HTTP server ready to accept connections
func NewHTTPServer(c *conf.Server, logger log.Logger) *khttp.Server {
	// Configure CORS with security best practices
	// Security: CORS configuration must follow browser security rules:
	// - If AllowedOrigins contains "*", AllowCredentials must be false
	// - If AllowCredentials is true, AllowedOrigins must specify exact origins (not "*")
	var corsOpts []handlers.CORSOption
	allowedOrigins := []string{"*"} // Default: allow all origins
	allowCredentials := false       // Default: disabled for security when using "*"

	// Security: If using "*" origin, credentials must be disabled
	// This prevents CORS security vulnerabilities
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
			// Note: Accept, Origin, and Access-Control-Request-* headers are simple headers
			// and don't need to be explicitly allowed
		}),
		handlers.MaxAge(86400), // Cache preflight requests for 24 hours
	)

	if allowCredentials {
		corsOpts = append(corsOpts, handlers.AllowCredentials())
	}

	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
		),
	}

	// Create CORS middleware factory (reusable)
	// This creates a single CORS middleware function that can wrap any handler
	corsMiddleware := handlers.CORS(corsOpts...)

	// Add CORS filter
	// This ensures CORS headers are applied to all responses
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
