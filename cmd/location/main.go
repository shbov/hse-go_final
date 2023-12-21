package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/location/app"
	"github.com/shbov/hse-go_final/internal/location/logger"
	"log"
)

func getEnvPath() string {
	var envPath string

	flag.StringVar(&envPath, "env", ".env", "path to .env file")
	flag.Parse()

	return envPath
}

func main() {
	envPath := getEnvPath()
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)

	}

	config, err := app.ParseConfigFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	lg, err := logger.GetLogger(config.App.Debug, config.App.DSN, "development")
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := zapctx.WithLogger(context.Background(), lg)
	a, err := app.New(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := a.Serve(ctx); err != nil {
		lg.Fatal(err.Error())
	}
}
