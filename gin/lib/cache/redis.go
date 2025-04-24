package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrLockNotAcquired 获取分布式锁失败
var ErrLockNotAcquired = errors.New("redis: lock not acquired")

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 单例管理
var (
	mu        sync.Mutex
	instances = make(map[string]*redis.Client)
)

// GetRedisClient 获取或创建单例
func GetRedisClient(name string, cfg *RedisConfig) *redis.Client {
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
	instances[name] = rdb
	return rdb
}
