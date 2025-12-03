package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/example/lib/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetDBMiddleware creates a middleware that injects a GORM database instance into the request context.
// The database can be retrieved later using factory.DB(ctx).
//
// Parameters:
//   - db: The GORM database instance to inject into the context
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware function that adds database to context
func SetDBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Request = req.WithContext(context.WithValue(req.Context(), utils.ContextKey("DB"), db))
		ctx.Next()
	}
}

// SetRedisMiddleware creates a middleware that injects a Redis client into the request context.
// The Redis client can be retrieved later using factory.Redis(ctx).
//
// Parameters:
//   - redis: The Redis client instance to inject into the context
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware function that adds Redis client to context
func SetRedisMiddleware(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Request = req.WithContext(context.WithValue(req.Context(), utils.ContextKey("REDIS"), redis))
		ctx.Next()
	}
}
