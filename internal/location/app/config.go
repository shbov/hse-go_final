package app

import (
	"github.com/shbov/hse-go_final/internal/location/httpadapter"
	"os"
	"time"
)

const (
	AppName                = "location"
	DefaultServeAddress    = "localhost:8080"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/location/v1"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/location"

	DefaultOtlpAddress = "localhost:4317"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	DSN             string        `yaml:"dsn"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DSN           string `yaml:"dsn"`
	MigrationsDir string `yaml:"migrations_dir"`
}

type Config struct {
	App      AppConfig          `yaml:"app"`
	Database DatabaseConfig     `yaml:"database"`
	HTTP     httpadapter.Config `yaml:"http"`
}

func ParseConfigFromEnv() (*Config, error) {
	return &Config{
		App: AppConfig{
			Debug:           getDebug(os.Getenv("APP_DEBUG")),
			DSN:             getEnv(os.Getenv("APP_DSN"), DefaultDSN),
			ShutdownTimeout: parseDuration(os.Getenv("APP_SHUTDOWN_TIMEOUT"), DefaultShutdownTimeout),
		},
		Database: DatabaseConfig{
			DSN:           getEnv(os.Getenv("LOCATION_DB_DSN"), DefaultDSN),
			MigrationsDir: getEnv(os.Getenv("LOCATION_DB_MIGRATIONS_DIR"), DefaultMigrationsDir),
		},
		HTTP: httpadapter.Config{
			ServeAddress:   getEnv(os.Getenv("LOCATION_HTTP_SERVE_ADDRESS"), DefaultServeAddress),
			BasePath:       getEnv(os.Getenv("LOCATION_HTTP_BASE_PATH"), DefaultBasePath),
			OtlpAddress:    getEnv(os.Getenv("HTTP_OTLP"), DefaultOtlpAddress),
			SwaggerAddress: getEnv(os.Getenv("LOCATION_HTTP_SWAGGER_ADDRESS"), ""),
		},
	}, nil
}

func getDebug(getenv string) bool {
	if getenv == "" {
		return false
	}

	return getenv == "true"
}

func getEnv(getenv string, address string) string {
	if getenv == "" {
		return address
	}

	return getenv
}

func parseDuration(getenv string, timeout time.Duration) time.Duration {
	if getenv == "" {
		return timeout
	}

	d, err := time.ParseDuration(getenv)
	if err != nil {
		return timeout
	}

	return d
}
