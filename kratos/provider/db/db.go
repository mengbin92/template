package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"explorer/internal/conf"
	"explorer/provider/db/mysql"
	"explorer/provider/db/postgres"
	"explorer/provider/db/sqlite3"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var (
	gdb        *gorm.DB
	initDBOnce sync.Once
)

// Init inits the database connection only once
func Init(ctx context.Context, cfg *conf.Database, logger log.Logger) error {
	var err error

	initDBOnce.Do(func() {
		if cfg.Driver == "postgre" {
			gdb, err = postgres.InitDB(cfg.Source)
		} else if cfg.Driver == "sqlite" {
			gdb, err = sqlite3.InitDB(cfg.Source)
		} else {
			gdb, err = mysql.InitDB(cfg.Source) // MySQL is default
		}
	})

	sqlDB, err := gdb.DB()
	if err != nil {
		panic(fmt.Errorf("set connection error: %s", err.Error()))
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return err
}

func Get() *gorm.DB {
	if gdb == nil {
		panic("db is nil")
	}

	return gdb
}
