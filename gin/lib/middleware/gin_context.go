package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/example/lib/utils"
	"go.uber.org/zap"
)

// SetLogMiddleware creates a middleware that injects a logger into the request context.
// The logger can be retrieved later using factory.Logger(ctx).
//
// Parameters:
//   - logger: The zap logger instance to inject into the context
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware function that adds logger to context
func SetLogMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Request = req.WithContext(context.WithValue(req.Context(), utils.ContextKey("LOGGER"), logger))
		ctx.Next()
	}
}
