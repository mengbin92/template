package cache

import (
	"context"
	"explorer/internal/conf"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

var (
	redisc        *redis.Client
	initRedisOnce sync.Once
)

func InitRedis(ctx context.Context, cfg *conf.Redis, logger log.Logger) error {
	initRedisOnce.Do(func() {
		redisc = redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           int(cfg.Db),
			PoolSize:     int(cfg.PoolSize),
			ReadTimeout:  cfg.ReadTimeout.AsDuration(),
			WriteTimeout: cfg.WriteTimeout.AsDuration(),
		})
	})
	logger.Log(log.LevelInfo, "msg", "init redis client ....")
	_, err := redisc.Ping(ctx).Result()
	if err != nil {
		return errors.Wrap(err, "redis ping error")
	}
	return nil
}

func GetRedisClient() *redis.Client {
	return redisc
}
