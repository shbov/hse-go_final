package app

import (
	"github.com/shbov/hse-go_final/internal/driver/httpadapter"
	"github.com/shbov/hse-go_final/pkg/config"
	"os"
	"time"
)

const (
	AppName                = "driver"
	DefaultServeAddress    = "localhost:8081"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/driver/v1/"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/location"

	DefaultOtlpAddress    = "localhost:4317"
	DefaultSwaggerAddress = "localhost:8081"
)

type AppConfig struct {
	Debug           bool          `yaml:"debug"`
	DSN             string        `yaml:"dsn"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type MongoConfig struct {
	Database string `yaml:"database"`
	Uri      string `yaml:"uri"`
}

type KafkaConfig struct {
	Brokers  []string `yaml:"brokers"`
	GroupID  string   `yaml:"group_id"`
	Topic    string   `yaml:"topic"`
	MaxBytes string   `yaml:"max_bytes"`
}

type Config struct {
	App  AppConfig          `yaml:"app"`
	HTTP httpadapter.Config `yaml:"http"`

	Environment string           `yaml:"environment"`
	ServiceName string           `yaml:"service_name"`
	Migrations  MigrationsConfig `yaml:"migration"`

	Mongo MongoConfig `yaml:"mongo"`
	Kafka KafkaConfig `yaml:"kafka"`
}

type MigrationsConfig struct {
	URI     string `yaml:"uri"`
	Path    string `yaml:"path"`
	Enabled bool   `yaml:"enabled"`
}

func ParseConfigFromEnv() (*Config, error) {
	return &Config{
		App: AppConfig{
			Debug:           config.GetDebug(os.Getenv("APP_DEBUG")),
			DSN:             config.GetEnv(os.Getenv("APP_DSN"), "dsn://"),
			ShutdownTimeout: config.ParseDuration(os.Getenv("APP_SHUTDOWN_TIMEOUT"), 20*time.Second),
		},

		HTTP: httpadapter.Config{
			ServeAddress:   config.GetEnv(os.Getenv("DRIVER_HTTP_SERVE_ADDRESS"), DefaultServeAddress),
			BasePath:       config.GetEnv(os.Getenv("DRIVER_HTTP_BASE_PATH"), DefaultBasePath),
			OtlpAddress:    config.GetEnv(os.Getenv("HTTP_OTLP"), DefaultOtlpAddress),
			SwaggerAddress: config.GetEnv(os.Getenv("DRIVER_HTTP_SWAGGER_ADDRESS"), DefaultSwaggerAddress),
		},

		Mongo: MongoConfig{
			Database: "driver",
			Uri:      "mongodb://localhost:27017",
		},

		Kafka: KafkaConfig{
			Brokers:  []string{"localhost:9092"},
			GroupID:  "driver",
			Topic:    "driver",
			MaxBytes: "2000000",
		},
	}, nil
}
