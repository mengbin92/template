// Package factory provides helper functions to retrieve resources from context.
// It extracts database, Redis, and logger instances that were injected via middleware.
package factory

import (
	"context"
	"fmt"

	"github.com/mengbin92/example/lib/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DB retrieves the GORM database instance from the context.
// The database should be injected into the context using middleware.SetDBMiddleware.
//
// Parameters:
//   - ctx: The context containing the database instance
//
// Returns:
//   - *gorm.DB: The GORM database instance
//
// Panics:
//   - If the database is not found in the context
//   - If the value in context is not a *gorm.DB
func DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(utils.ContextKey("DB"))
	if v == nil {
		panic("db is not exist in context")
	}

	db, ok := v.(*gorm.DB)
	if !ok {
		panic(fmt.Sprintf("db in context is not *gorm.DB, got %T", v))
	}

	return db
}

// Redis retrieves the Redis client instance from the context.
// The Redis client should be injected into the context using middleware.SetRedisMiddleware.
//
// Parameters:
//   - ctx: The context containing the Redis client instance
//
// Returns:
//   - *redis.Client: The Redis client instance
//
// Panics:
//   - If the Redis client is not found in the context
//   - If the value in context is not a *redis.Client
func Redis(ctx context.Context) *redis.Client {
	v := ctx.Value(utils.ContextKey("REDIS"))
	if v == nil {
		panic("redis is not exist in context")
	}

	redisClient, ok := v.(*redis.Client)
	if !ok {
		panic(fmt.Sprintf("redis in context is not *redis.Client, got %T", v))
	}

	return redisClient
}

// Logger retrieves the zap logger instance from the context.
// The logger should be injected into the context using middleware.SetLogMiddleware.
//
// Parameters:
//   - ctx: The context containing the logger instance
//
// Returns:
//   - *zap.Logger: The zap logger instance
//
// Panics:
//   - If the logger is not found in the context
//   - If the value in context is not a *zap.Logger
func Logger(ctx context.Context) *zap.Logger {
	v := ctx.Value(utils.ContextKey("LOGGER"))
	if v == nil {
		panic("zap.Logger is not exist in context")
	}

	log, ok := v.(*zap.Logger)
	if !ok {
		panic(fmt.Sprintf("logger in context is not *zap.Logger, got %T", v))
	}

	return log
}

// DBOrNil retrieves the GORM database instance from the context, returning nil if not found.
// This is a safe version of DB that does not panic.
//
// Parameters:
//   - ctx: The context containing the database instance
//
// Returns:
//   - *gorm.DB: The GORM database instance, or nil if not found
func DBOrNil(ctx context.Context) *gorm.DB {
	v := ctx.Value(utils.ContextKey("DB"))
	if v == nil {
		return nil
	}

	db, ok := v.(*gorm.DB)
	if !ok {
		return nil
	}

	return db
}

// RedisOrNil retrieves the Redis client instance from the context, returning nil if not found.
// This is a safe version of Redis that does not panic.
//
// Parameters:
//   - ctx: The context containing the Redis client instance
//
// Returns:
//   - *redis.Client: The Redis client instance, or nil if not found
func RedisOrNil(ctx context.Context) *redis.Client {
	v := ctx.Value(utils.ContextKey("REDIS"))
	if v == nil {
		return nil
	}

	redisClient, ok := v.(*redis.Client)
	if !ok {
		return nil
	}

	return redisClient
}

// LoggerOrNil retrieves the zap logger instance from the context, returning nil if not found.
// This is a safe version of Logger that does not panic.
//
// Parameters:
//   - ctx: The context containing the logger instance
//
// Returns:
//   - *zap.Logger: The zap logger instance, or nil if not found
func LoggerOrNil(ctx context.Context) *zap.Logger {
	v := ctx.Value(utils.ContextKey("LOGGER"))
	if v == nil {
		return nil
	}

	log, ok := v.(*zap.Logger)
	if !ok {
		return nil
	}

	return log
}
