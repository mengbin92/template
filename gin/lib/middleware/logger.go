package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        // 处理请求
        c.Next()

        latency := time.Since(start)
        logger.Info("HTTP Request",
            zap.Int("status", c.Writer.Status()),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.String("ip", c.ClientIP()),
            zap.String("user-agent", c.Request.UserAgent()),
            zap.Duration("latency", latency),
            zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
        )
    }
}
