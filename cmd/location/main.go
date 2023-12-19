package main

import (
	"context"
	"flag"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/location/app"
	"github.com/shbov/hse-go_final/internal/location/logger"
	"log"
)

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "./.config/location/config.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	path := getConfigPath()
	config, err := app.NewConfig(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	lg, err := logger.GetLogger(true, config.App.DSN, "development")
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := zapctx.WithLogger(context.Background(), lg)
	_, err = app.New(ctx, config)
	if err != nil {
		log.Fatal(err.Error())
	}
}
