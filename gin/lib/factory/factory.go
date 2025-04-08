package factory

import (
	"context"

	"github.com/mengbin92/example/lib/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(utils.ContextKey("DB"))
	if v == nil {
		panic("db is not exist")
	}
	if db, ok := v.(*gorm.DB); ok {
		return db
	}
	panic("db is not exist")
}

func Redis(ctx context.Context) *redis.Client {
	v := ctx.Value(utils.ContextKey("REDIS"))
	if v == nil {
		panic("redis is not exist")
	}
	if redis, ok := v.(*redis.Client); ok {
		return redis
	}
	panic("redis is not exist")
}

func Logger(ctx context.Context) *zap.Logger {
	v := ctx.Value(utils.ContextKey("LOGGER"))
	if v == nil {
		panic("zap.Logger is not exist")
	}
	if log, ok := v.(*zap.Logger); ok {
		return log
	}
	panic("zap.Logger is not exist")
}
