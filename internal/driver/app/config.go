package app

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"time"
)

const (
	EnvProduction = "PRODUCTION"
	EnvPrefix     = "DRIVER"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type MongoConfig struct {
	Database string `yaml:"database"          env:"DEMO_MONGO_DATABASE" default:"demo"`
	Uri      string `yaml:"uri"               env:"DEMO_MONGO_URI"`
}

type Config struct {
	Environment string           `yaml:"environment" env:"DRIVER_ENVIRONMENT" default:"development"`
	ServiceName string           `yaml:"service_name" env:"DRIVER_SERVICE_NAME" default:"driver-server"`
	Server      ServerConfig     `yaml:"server"         env:"-"`
	Mongo       MongoConfig      `yaml:"mongo" env:"-"`
	Migrations  MigrationsConfig `yaml:"migration" env:"-"`
}

type ServerConfig struct {
	HttpServerPort int `yaml:"http_server_port" env:"DRIVER_HTTP_SERVER_PORT" default:"8080"`
}

type MigrationsConfig struct {
	URI     string `yaml:"uri"     env:"DRIVER_MONGO_MIGRATION_URI"`
	Path    string `yaml:"path"    env:"DRIVER_MIGRATIONS_PATH"`
	Enabled bool   `yaml:"enabled" env:"DRIVER_MIGRATIONS_ENABLED"   default:"false"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config
	loader := configor.New(&configor.Config{ENVPrefix: EnvPrefix, ErrorOnUnmatchedKeys: true})
	if err := loader.Load(&cfg, path); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Environment != EnvProduction {
		spew.Dump(cfg)
	}

	return &cfg, nil
}
