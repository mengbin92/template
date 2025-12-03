// Package global provides global variables and initialization functions for the application.
// It manages the logger and other global resources that are used throughout the application.
package global

import (
	"context"

	"kratos-project-template/internal/conf"
	"kratos-project-template/provider/cache"
	"kratos-project-template/provider/db"
	"kratos-project-template/provider/storage"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	// Logger is the global logger instance used throughout the application.
	Logger *log.Helper
)

// Init initializes global variables including the logger.
// It initializes database, cache, and object storage connections based on the bootstrap configuration.
//
// Parameters:
//   - bc: The bootstrap configuration containing log, data, and other settings
//   - logger: The logger instance to use for application logging
//
// The function will panic if critical initialization steps fail:
//   - Bootstrap configuration is nil
//   - Database initialization fails
func Init(bc *conf.Bootstrap, logger log.Logger) {
	if bc == nil {
		panic("bootstrap config cannot be nil")
	}

	Logger = log.NewHelper(logger)
	Logger.Infof("logger initialized: %v", bc.Log)

	Logger.Infof("database initialized")
	err := db.Init(context.Background(), bc.Data.Database, logger)
	if err != nil {
		panic(err)
	}

	err = cache.InitRedis(context.Background(), bc.Data.Redis, logger)
	if err != nil {
		Logger.Warnf("redis initialization failed: %v", err)
	}

	Logger.Infof("object storage initialized")
	if bc.Data != nil {
		err = storage.Init(context.Background(), bc.Data.GetObjectStorage(), logger)
		if err != nil {
			Logger.Warnf("object storage initialization failed: %v", err)
		}
	}
}

