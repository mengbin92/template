package main

import (
	"context"
	"explorer/internal/conf"
	"explorer/internal/models/users"
	"explorer/internal/server"
	"explorer/internal/service"
	"explorer/provider/cache"
	"explorer/provider/db"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"

	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(ctx context.Context, bc *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	// init db
	if err := db.Init(ctx, bc.Database, logger); err != nil {
		logger.Log(log.LevelFatal, "msg", "init db failed")
		return nil, nil, errors.Wrap(err, "init db failed")
	}
	db.Get().AutoMigrate(&users.User{})

	// init redis
	if err := cache.InitRedis(ctx, bc.Redis, logger); err != nil {
		logger.Log(log.LevelFatal, "msg", "init redis failed")
		return nil, nil, errors.Wrap(err, "init redis failed")
	}

	cleanup := func(ctx context.Context) {
		logger.Log(log.LevelInfo, "msg", "close the data resources")

		// close redis
		cache.GetRedisClient().Close()
	}

	// init service
	basicClient := service.NewBasicService(logger)

	userManager := users.NewUserManager(bc.AuthConfig, logger)
	userService := service.NewUserService(userManager, logger)

	httpServer := server.NewHTTPServer(bc.Server, basicClient, userService, logger)
	grpcServer := server.NewGRPCServer(bc.Server, basicClient, userService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup(ctx)
	}, nil
}
