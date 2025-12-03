// Package middleware provides HTTP middleware functions for the Gin framework.
// It includes logging, database context injection, and other request processing utilities.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetLoggerMiddleware creates a middleware that logs HTTP request details.
// It records request method, path, query parameters, client IP, user agent, latency, and errors.
//
// Parameters:
//   - logger: The zap logger instance to use for logging
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware function that logs request information
func SetLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
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
