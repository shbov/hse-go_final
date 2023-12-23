package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/app"
	"github.com/shbov/hse-go_final/internal/driver/config"
	"github.com/shbov/hse-go_final/pkg/logger"
	"log"
)

func getEnvPath() string {
	var envPath string

	flag.StringVar(&envPath, "env", ".env", "path to .env file")
	flag.Parse()

	return envPath
}

// TODO: проверить тесты, поднять % покрытия
// TODO: проверить весь путь событий и команд по схеме: создание поездки -> ...
// TODO: проверить, что в сваггере указаны правильные методы и отписания
// + TODO: сделать long polling

// TODO: наполнить БД больше тестовыми данными, проверить корректнось всех запросов (корнер-кейсы в частности)
func main() {
	envPath := getEnvPath()
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	var env string
	if cfg.App.Debug {
		env = "development"
	} else {
		env = "production"
	}

	lg, err := logger.GetLogger(cfg.App.Debug, cfg.App.DSN, env, cfg.App.AppName)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := zapctx.WithLogger(context.Background(), lg)
	a, err := app.New(ctx, cfg)

	if err != nil {
		lg.Fatal(err.Error())
	}

	if err := a.Serve(ctx); err != nil {
		lg.Error(err.Error())
	}
}
