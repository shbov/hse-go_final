package app

import (
	"github.com/shbov/hse-go_final/internal/location/httpadapter"
	"github.com/shbov/hse-go_final/pkg/config"
	"os"
	"time"
)

const (
	DefaultAppName         = "location"
	DefaultServeAddress    = "localhost:8080"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/location/v1"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/location"

	DefaultOtlpAddress    = "localhost:4317"
	DefaultSwaggerAddress = "localhost:8080"
)

type AppConfig struct {
	AppName         string        `yaml:"app_name"`
	Debug           bool          `yaml:"debug"`
	DSN             string        `yaml:"dsn"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
	Migrate       bool   `yaml:"migrate"`
	Populate      bool   `yaml:"populate"`
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`
}

func ParseConfigFromEnv() (*Config, error) {
	return &Config{
		App: AppConfig{
			AppName:         DefaultAppName,
			Debug:           config.GetEnvBoolean(os.Getenv("APP_DEBUG")),
			DSN:             config.GetEnv(os.Getenv("APP_DSN"), DefaultDSN),
			ShutdownTimeout: config.ParseDuration(os.Getenv("APP_SHUTDOWN_TIMEOUT"), DefaultShutdownTimeout),
		},
		Database: DatabaseConfig{
			DSN:           config.GetEnv(os.Getenv("LOCATION_DB_DSN"), DefaultDSN),
			MigrationsDir: config.GetEnv(os.Getenv("LOCATION_DB_MIGRATIONS_DIR"), DefaultMigrationsDir),
			Migrate:       config.GetEnvBoolean(os.Getenv("LOCATION_DB_MIGRATE")),
			Populate:      config.GetEnvBoolean(os.Getenv("LOCATION_DB_POPULATE")),
		},
		HTTP: httpadapter.Config{
			ServeAddress:   config.GetEnv(os.Getenv("LOCATION_HTTP_SERVE_ADDRESS"), DefaultServeAddress),
			BasePath:       config.GetEnv(os.Getenv("LOCATION_HTTP_BASE_PATH"), DefaultBasePath),
			OtlpAddress:    config.GetEnv(os.Getenv("HTTP_OTLP"), DefaultOtlpAddress),
			SwaggerAddress: config.GetEnv(os.Getenv("LOCATION_HTTP_SWAGGER_ADDRESS"), DefaultSwaggerAddress),
		},
	}, nil
}
