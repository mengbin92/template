// Package db provides database connection initialization and management.
// It supports multiple database drivers including MySQL, PostgreSQL, and SQLite.
package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/mengbin92/example/lib/db/mysql"
	"github.com/mengbin92/example/lib/db/postgres"
	"github.com/mengbin92/example/lib/db/sqlite3"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// gdb is the global GORM database instance
	gdb *gorm.DB
	// initDBOnce ensures the database is initialized only once (thread-safe)
	initDBOnce sync.Once
)

// Init initializes the database connection.
// It uses sync.Once to ensure the database is initialized only once, even if called multiple times.
//
// Supported drivers:
//   - "postgre": PostgreSQL database
//   - "sqlite": SQLite database
//   - Default: MySQL database
//
// Parameters:
//   - driver: The database driver type ("postgre", "sqlite", or empty/default for MySQL)
//   - source: The database connection source string (DSN)
//
// Returns:
//   - error: Error if initialization or connection fails
func Init(driver, source string) error {
	if source == "" {
		return errors.New("database source cannot be empty")
	}

	var initErr error
	gormLogger := NewGormLogger(logger.Error)

	initDBOnce.Do(func() {
		// Retry connection with exponential backoff
		maxRetries := 10
		retryDelay := 2 * time.Second

		for attempt := 0; attempt < maxRetries; attempt++ {
			switch driver {
			case "postgre":
				gdb, initErr = postgres.InitDB(source, gormLogger)
			case "sqlite":
				gdb, initErr = sqlite3.InitDB(source, gormLogger)
			default:
				// MySQL is the default driver
				gdb, initErr = mysql.InitDB(source, gormLogger)
			}

			if initErr == nil {
				// Connection successful
				break
			}

			// Log retry attempt
			if attempt < maxRetries-1 {
				fmt.Printf("database connection attempt %d/%d failed: %v, retrying in %v...\n",
					attempt+1, maxRetries, initErr, retryDelay)
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
			}
		}
	})

	if initErr != nil {
		return errors.Wrap(initErr, "connect to db error")
	}

	if gdb == nil {
		return errors.New("database instance is nil after initialization")
	}

	// Configure connection pool settings
	sqlDB, err := gdb.DB()
	if err != nil {
		return errors.Wrap(err, "get sql db error")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping database error")
	}

	return nil
}

// Get returns the global GORM database instance.
//
// Returns:
//   - *gorm.DB: The GORM database instance
//
// Panics:
//   - If the database has not been initialized (nil instance)
//
// Note: The database should be initialized using Init before calling this function.
func Get() *gorm.DB {
	if gdb == nil {
		panic("database is nil; please initialize it using db.Init()")
	}

	return gdb
}

// NewGormLogger creates a GORM logger adapter.
// This creates a simple GORM logger using the standard log package.
//
// Parameters:
//   - level: The GORM log level
//
// Returns:
//   - logger.Interface: A GORM logger interface implementation
func NewGormLogger(level logger.LogLevel) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
			SlowThreshold:             time.Second,
		},
	)
}
