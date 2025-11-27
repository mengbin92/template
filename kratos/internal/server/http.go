package server

import (
	"encoding/json"
	v1 "explorer/api/explorer/v1"
	"explorer/internal/conf"
	"explorer/internal/middleware"
	"explorer/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v3"
)

var (
	openAPIFileCache     string
	openAPIFileCacheOnce sync.Once
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, basicService *service.BasicService, userService *service.UserService, logger log.Logger) *khttp.Server {
	// Create swagger handler first (will be used in filter)
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/openapi.json"), // Point to our OpenAPI JSON file
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	)

	// Configure CORS with security best practices
	// Note: If AllowCredentials is true, AllowedOrigins cannot contain "*"
	var corsOpts []handlers.CORSOption
	allowedOrigins := []string{"*"} // Default: allow all origins
	allowCredentials := false       // Default: don't allow credentials for security

	// If credentials are needed, origins must be explicitly specified (not "*")
	// This prevents CORS security vulnerabilities
	if c.Http != nil {
		// For now, use default CORS settings
		// In production, consider adding CORS configuration to Server_HTTP proto
		// and reading from config file
	}

	// Security: If using "*" origin, credentials must be disabled
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		allowCredentials = false
	}

	corsOpts = append(corsOpts,
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
	)

	if allowCredentials {
		corsOpts = append(corsOpts, handlers.AllowCredentials())
	}

	authMiddleware := middleware.NewAuthMiddleware(userService.UserManager, logger)

	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			middleware.AuthMiddlewareWrap(authMiddleware),
		),
		khttp.Filter(handlers.CORS(corsOpts...)),
		// Add filter to handle swagger requests before routing
		khttp.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				path := r.URL.Path
				// Handle swagger requests
				if path == "/swagger" {
					http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
					return
				}
				if strings.HasPrefix(path, "/swagger/") || path == "/swagger" {
					swaggerHandler.ServeHTTP(w, r)
					return
				}
				// Continue with normal routing
				next.ServeHTTP(w, r)
			})
		}),
	}
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
	v1.RegisterBasicHTTPServer(srv, basicService)
	v1.RegisterUserHTTPServer(srv, userService)

	// Register Swagger UI routes
	registerSwaggerRoutes(srv)

	return srv
}

// registerSwaggerRoutes registers Swagger UI routes for OpenAPI documentation
func registerSwaggerRoutes(srv *khttp.Server) {
	r := srv.Route("/")

	// Serve OpenAPI YAML file
	r.GET("/openapi.yaml", func(ctx khttp.Context) error {
		openapiPath, err := findOpenAPIFile()
		if err != nil {
			ctx.Response().WriteHeader(404)
			ctx.Response().Write([]byte("openapi.yaml not found"))
			return nil
		}

		data, err := os.ReadFile(openapiPath)
		if err != nil {
			ctx.Response().WriteHeader(500)
			ctx.Response().Write([]byte("Failed to read openapi.yaml: " + err.Error()))
			return nil
		}

		ctx.Response().Header().Set("Content-Type", "application/x-yaml")
		ctx.Response().Write(data)
		return nil
	})

	// Serve OpenAPI JSON file (converted from YAML)
	r.GET("/openapi.json", func(ctx khttp.Context) error {
		openapiPath, err := findOpenAPIFile()
		if err != nil {
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusNotFound)
			errorMsg := map[string]string{
				"error":   "openapi.yaml not found",
				"message": "The openapi.yaml file could not be located. Please ensure it exists in the project root directory.",
			}
			if jsonData, err := json.Marshal(errorMsg); err == nil {
				ctx.Response().Write(jsonData)
			} else {
				ctx.Response().Write([]byte(`{"error":"openapi.yaml not found"}`))
			}
			return nil
		}

		// Security: Validate that the path is within project root to prevent path traversal
		projectRoot := findProjectRoot()
		cleanPath := filepath.Clean(openapiPath)
		cleanRoot := filepath.Clean(projectRoot)

		// Ensure the resolved path is within project root
		relPath, err := filepath.Rel(cleanRoot, cleanPath)
		if err != nil || strings.HasPrefix(relPath, "..") {
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusForbidden)
			ctx.Response().Write([]byte(`{"error":"invalid file path"}`))
			return nil
		}

		data, err := os.ReadFile(cleanPath)
		if err != nil {
			// Don't leak internal path information
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			errorMsg := map[string]string{
				"error":   "failed to read openapi.yaml",
				"message": "Unable to read the openapi.yaml file",
			}
			if jsonData, err := json.Marshal(errorMsg); err == nil {
				ctx.Response().Write(jsonData)
			} else {
				ctx.Response().Write([]byte(`{"error":"failed to read openapi.yaml"}`))
			}
			return nil
		}

		// Parse YAML
		var yamlData interface{}
		if err := yaml.Unmarshal(data, &yamlData); err != nil {
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			errorMsg := map[string]string{
				"error":   "failed to parse openapi.yaml",
				"message": "Unable to parse the openapi.yaml file",
			}
			if jsonData, err := json.Marshal(errorMsg); err == nil {
				ctx.Response().Write(jsonData)
			} else {
				ctx.Response().Write([]byte(`{"error":"failed to parse openapi.yaml"}`))
			}
			return nil
		}

		// Convert to JSON
		jsonData, err := json.Marshal(yamlData)
		if err != nil {
			ctx.Response().Header().Set("Content-Type", "application/json")
			ctx.Response().WriteHeader(http.StatusInternalServerError)
			errorMsg := map[string]string{
				"error":   "failed to convert to JSON",
				"message": "Unable to convert openapi.yaml to JSON format",
			}
			if jsonData, err := json.Marshal(errorMsg); err == nil {
				ctx.Response().Write(jsonData)
			} else {
				ctx.Response().Write([]byte(`{"error":"failed to convert to JSON"}`))
			}
			return nil
		}

		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.Response().Write(jsonData)
		return nil
	})

	// Swagger routes are now handled by the HTTP Filter in NewHTTPServer
	// No need to register them here as the filter intercepts all /swagger/* requests
}

// findOpenAPIFile tries to locate the openapi.yaml file in common locations
func findOpenAPIFile() (string, error) {
	// Use cached result if available
	openAPIFileCacheOnce.Do(func() {
		projectRoot := findProjectRoot()
		openAPIPath := filepath.Join(projectRoot, "openapi.yaml")
		if info, err := os.Stat(openAPIPath); err == nil && !info.IsDir() {
			openAPIFileCache = openAPIPath
			return
		}

		// Try other common locations
		cwd, _ := os.Getwd()
		paths := []string{
			filepath.Join(projectRoot, "openapi.yaml"),
			filepath.Join(cwd, "openapi.yaml"),
			filepath.Join(cwd, "..", "openapi.yaml"),
			filepath.Join(cwd, "..", "..", "openapi.yaml"),
			"openapi.yaml",
		}

		if execPath, err := os.Executable(); err == nil {
			execDir := filepath.Dir(execPath)
			paths = append(paths,
				filepath.Join(execDir, "openapi.yaml"),
				filepath.Join(execDir, "..", "openapi.yaml"),
			)
		}

		for _, path := range paths {
			cleanPath := filepath.Clean(path)
			if info, err := os.Stat(cleanPath); err == nil && !info.IsDir() {
				openAPIFileCache = cleanPath
				return
			}
		}
	})

	if openAPIFileCache != "" {
		return openAPIFileCache, nil
	}

	return "", os.ErrNotExist
}

// findProjectRoot tries to find the project root by looking for go.mod file
func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Start from current directory and walk up to find go.mod
	dir := cwd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if info, err := os.Stat(goModPath); err == nil && !info.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, stop
			break
		}
		dir = parent
	}

	// Fallback: try executable directory
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		goModPath := filepath.Join(execDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return execDir
		}
		// Try parent directories
		dir := execDir
		for i := 0; i < 5; i++ {
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			goModPath := filepath.Join(parent, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				return parent
			}
			dir = parent
		}
	}

	return cwd // Fallback to current working directory
}
