package cmd

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"github.com/mengbin92/example/config"
	"github.com/mengbin92/example/lib/db"
	"github.com/mengbin92/example/lib/logger"
	"github.com/mengbin92/example/lib/middleware"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Execute() {
	// TODO: implement

	// 初始化配置
	config.LoadConfig()

	run()
}

func run() {
	if err := setEngine(loadDB()).Run(fmt.Sprintf(":%d", viper.GetInt("server.port"))); err != nil {
		log.Error("Failed to run server: ", err)
		panic(err)
	}

}

func loadDB() *gorm.DB {
	err := db.Init(viper.GetString("database.driver"), viper.GetString("database.source"))
	if err != nil {
		log.Error("Failed to connect to database: ", err)
		panic(err)
	}
	return db.Get()
}

func setEngine(db *gorm.DB) *gin.Engine {
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.New()

	// 注册自定义模板函数
	r.SetFuncMap(template.FuncMap{
		"formatUnixTime": func(ts string) string {
			timestamp, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				return "时间格式错误"
			}
			return formatUnixTime(timestamp)
		},
	})

	// 设置中间件
	r.Use(gin.Recovery())
	r.Use(middleware.SetLoggerMiddleware(logger.DefaultLogger(viper.GetInt("log.level"), viper.GetString("log.format"))))
	r.Use(middleware.SetDBMiddleware(db))
	r.Use(middleware.SetLogMiddleware(logger.DefaultLogger(viper.GetInt("log.level"), viper.GetString("log.format"))))

	//设置路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func formatUnixTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}
