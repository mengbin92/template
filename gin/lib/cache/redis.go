// Package cache provides Redis cache client initialization and management.
package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

// ErrLockNotAcquired indicates that a distributed lock could not be acquired.
var ErrLockNotAcquired = errors.New("redis: lock not acquired")

// RedisConfig contains configuration for a Redis client connection.
type RedisConfig struct {
	// Addr is the Redis server address (e.g., "localhost:6379")
	Addr string
	// Password is the Redis server password (empty if no password)
	Password string
	// DB is the Redis database number (0-15)
	DB int
	// PoolSize is the maximum number of socket connections
	PoolSize int
	// DialTimeout is the timeout for establishing connections
	DialTimeout time.Duration
	// ReadTimeout is the timeout for socket reads
	ReadTimeout time.Duration
	// WriteTimeout is the timeout for socket writes
	WriteTimeout time.Duration
}

var (
	// mu protects the instances map from concurrent access
	mu sync.Mutex
	// instances stores Redis client instances by name (singleton pattern)
	instances = make(map[string]*redis.Client)
)

// GetRedisClient returns or creates a Redis client instance (singleton pattern).
// If a client with the given name already exists, it returns the existing instance.
// Otherwise, it creates a new client with the provided configuration.
//
// Parameters:
//   - name: The name identifier for the Redis client instance
//   - cfg: Redis configuration containing address, password, database, and pool settings
//
// Returns:
//   - *redis.Client: A configured Redis client instance
//
// Note: The function is thread-safe and uses a mutex to protect concurrent access.
func GetRedisClient(name string, cfg *RedisConfig) *redis.Client {
	if cfg == nil {
		panic("redis config cannot be nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if c, ok := instances[name]; ok {
		return c
	}

	opts := &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	rdb := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		// Log error but don't fail - connection might succeed later
		// In production, you might want to handle this differently
		_ = err
	}

	instances[name] = rdb
	return rdb
}

// InitRedis initializes a Redis client with the given name and configuration.
// It tests the connection and returns an error if the connection fails.
//
// Parameters:
//   - name: The name identifier for the Redis client instance
//   - cfg: Redis configuration
//
// Returns:
//   - error: Error if initialization or connection test fails
func InitRedis(name string, cfg *RedisConfig) error {
	if cfg == nil {
		return errors.New("redis config cannot be nil")
	}

	client := GetRedisClient(name, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return pkgerrors.Wrap(err, "redis ping error")
	}

	return nil
}
