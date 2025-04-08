package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mengbin92/example/lib/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetDBMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Request = req.WithContext(context.WithValue(req.Context(), utils.ContextKey("DB"), db))
		ctx.Next()
	}
}

func SetRedisMiddleware(redis *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Request = req.WithContext(context.WithValue(req.Context(), utils.ContextKey("REDIS"), redis))
		ctx.Next()
	}
}
