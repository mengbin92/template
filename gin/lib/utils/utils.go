// Package utils provides utility functions and types used throughout the application.
package utils

// ContextKey is a type-safe key for storing values in context.Context.
// It prevents key collisions when storing different types of values in the same context.
//
// Usage:
//   ctx := context.WithValue(ctx, ContextKey("DB"), dbInstance)
//   ctx := context.WithValue(ctx, ContextKey("REDIS"), redisClient)
//   ctx := context.WithValue(ctx, ContextKey("LOGGER"), logger)
type ContextKey string
