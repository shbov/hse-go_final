package main

import (
	"context"
	"flag"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/location/app"
	"github.com/shbov/hse-go_final/internal/location/logger"
	"log"
)

var DSN = ""

func getConfigPath() string {
	var configPath string

	flag.StringVar(&configPath, "c", "./.config/location/config.yaml", "path to config file")
	flag.Parse()

	return configPath
}

func main() {
	lg, err := logger.GetLogger(true, DSN, "development")
	if err != nil {
		log.Fatal(err.Error())
	}

	path := getConfigPath()
	config, err := app.NewConfig(path)
	if err != nil {
		lg.Fatal(err.Error())
	}

	ctx := zapctx.WithLogger(context.Background(), lg)
	_, err = app.New(ctx, config)
	if err != nil {
		lg.Fatal(err.Error())
	}
}
