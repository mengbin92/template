package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/mengbin92/example/lib/db/mysql"
	"github.com/mengbin92/example/lib/db/postgres"
	"github.com/mengbin92/example/lib/db/sqlite3"
	"gorm.io/gorm"
)

var (
	gdb      *gorm.DB
	initDBOnce sync.Once
)

// Init inits the database connection only once
func Init(driver, source string) error {
	var err error

	initDBOnce.Do(func() {
		if driver == "postgre" {
			gdb, err = postgres.InitDB(source)
		} else if driver == "sqlite" {
			gdb, err = sqlite3.InitDB(source)
		} else {
			gdb, err = mysql.InitDB(source) // MySQL is default
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
