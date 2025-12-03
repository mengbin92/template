// Package main is the entry point of the application.
// It initializes configuration, logging, and starts the gRPC and HTTP servers.
package main

import (
	"flag"
	"fmt"
	"os"

	"kratos-project-template/internal/conf"
	"kratos-project-template/internal/global"
	"kratos-project-template/provider/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Build-time variables set via ldflags.
var (
	Name     string
	Version  string
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *khttp.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource(""),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to scan configuration: %v\n", err)
		os.Exit(1)
	}

	logger := logger.NewZapLogger(bc.Log)

	logger = logger.With(
		zap.String("service.id", id),
		zap.String("service.name", Name),
		zap.String("service.version", Version),
	)

	// Initialize global variables
	global.Init(&bc, logger)

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		log.NewHelper(logger).Errorf("Failed to wire application: %v", err)
		os.Exit(1)
	}
	defer cleanup()

	// Start and wait for stop signal
	if err := app.Run(); err != nil {
		log.NewHelper(logger).Errorf("Application run error: %v", err)
		os.Exit(1)
	}
}
