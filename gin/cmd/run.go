// Package cmd provides command-line interface for the application.
// It handles server initialization, configuration loading, and HTTP server setup.
package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/mengbin92/example/config"
	"github.com/mengbin92/example/lib/cache"
	"github.com/mengbin92/example/lib/db"
	"github.com/mengbin92/example/lib/logger"
	"github.com/mengbin92/example/lib/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Execute is the main entry point for the command-line interface.
// It loads configuration and starts the HTTP server.
//
// The function will panic if critical initialization steps fail:
//   - Configuration loading fails
//   - Server startup fails
func Execute() {
	// Load configuration
	config.LoadConfig()

	run()
}

// run initializes the database connection and starts the HTTP server.
//
// The function will panic if:
//   - Database initialization fails
//   - Server startup fails
func run() {
	dbInstance := loadDB()
	engine := setEngine(dbInstance)

	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	if err := engine.Run(addr); err != nil {
		log.Error("Failed to run server: ", err)
		fmt.Fprintf(os.Stderr, "Failed to run server: %v\n", err)
		os.Exit(1)
	}
}

// loadRedis initializes and returns a Redis client instance.
//
// Parameters:
//   - name: The name identifier for the Redis client instance
//   - cfg: Redis configuration containing address, password, database, and pool settings
//
// Returns:
//   - *redis.Client: A configured Redis client instance
func loadRedis(name string, cfg *cache.RedisConfig) *redis.Client {
	return cache.GetRedisClient(name, cfg)
}

// loadDB initializes the database connection and returns the GORM database instance.
//
// Returns:
//   - *gorm.DB: The GORM database instance
//
// The function will panic if:
//   - Database initialization fails
//   - Database connection fails
func loadDB() *gorm.DB {
	driver := viper.GetString("database.driver")
	source := viper.GetString("database.source")

	if err := db.Init(driver, source); err != nil {
		log.Error("Failed to connect to database: ", err)
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	return db.Get()
}

// setEngine creates and configures a new Gin HTTP server engine.
// It sets up middleware, custom template functions, and routes.
//
// Parameters:
//   - db: The GORM database instance to inject into the request context
//
// Returns:
//   - *gin.Engine: A configured Gin engine ready to handle HTTP requests
func setEngine(db *gorm.DB) *gin.Engine {
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.New()

	// Register custom template functions
	r.SetFuncMap(template.FuncMap{
		"formatUnixTime": func(ts string) string {
			timestamp, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				return "时间格式错误"
			}
			return formatUnixTime(timestamp)
		},
	})

	// Setup middleware
	r.Use(gin.Recovery())

	logLevel := viper.GetInt("log.level")
	logFormat := viper.GetString("log.format")
	zapLogger := logger.DefaultLogger(logLevel, logFormat)

	r.Use(middleware.SetLoggerMiddleware(zapLogger))
	r.Use(middleware.SetDBMiddleware(db))
	r.Use(middleware.SetLogMiddleware(zapLogger))

	// Setup routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

// formatUnixTime formats a Unix timestamp to a human-readable date-time string.
//
// Parameters:
//   - ts: Unix timestamp in seconds
//
// Returns:
//   - string: Formatted date-time string in "2006-01-02 15:04:05" format
func formatUnixTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
