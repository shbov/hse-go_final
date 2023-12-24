package config

import (
	"github.com/shbov/hse-go_final/internal/driver/httpadapter"
	"github.com/shbov/hse-go_final/pkg/config"
	"os"
	"time"
)

const (
	DefaultAppName         = "driver"
	DefaultServeAddress    = "localhost:8081"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/driver/v1/"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/driver"

	DefaultMigratePolicy  = true
	DefaultPopulatePolicy = true

	DefaultOtlpAddress    = "localhost:4317"
	DefaultSwaggerAddress = "localhost:8081"

	DefaultDatabase = "driver"
	DefaultMongoUri = "mongodb://localhost:27017"

	DefaultKafkaGroupID  = "driver"
	DefaultKafkaTopic    = "driver"
	DefaultKafkaMaxBytes = "2000000"

	DefaultLocationService = "http://localhost:8080"
)

var DefaultKafkaBrokers = []string{"kafka:29092"}

type AppConfig struct {
	AppName         string        `yaml:"app_name"`
	Debug           bool          `yaml:"debug"`
	DSN             string        `yaml:"dsn"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type MongoConfig struct {
	Database      string `yaml:"database"`
	Uri           string `yaml:"uri"`
	MigrationsDir string `yaml:"migrations_dir"`
	Migrate       bool   `yaml:"migrate"`
	Populate      bool   `yaml:"populate"`
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

	Mongo MongoConfig `yaml:"mongo"`
	Kafka KafkaConfig `yaml:"kafka"`

	Location LocationService `yaml:"location"`
}

type LocationService struct {
	URL string `yaml:"url"`
}

func ParseConfigFromEnv() (*Config, error) {
	return &Config{
		App: AppConfig{
			AppName:         DefaultAppName,
			Debug:           config.GetEnvBoolean(os.Getenv("APP_DEBUG")),
			DSN:             config.GetEnv(os.Getenv("APP_DSN"), DefaultDSN),
			ShutdownTimeout: config.ParseDuration(os.Getenv("APP_SHUTDOWN_TIMEOUT"), DefaultShutdownTimeout),
		},

		HTTP: httpadapter.Config{
			ServeAddress:   config.GetEnv(os.Getenv("DRIVER_HTTP_SERVE_ADDRESS"), DefaultServeAddress),
			BasePath:       config.GetEnv(os.Getenv("DRIVER_HTTP_BASE_PATH"), DefaultBasePath),
			OtlpAddress:    config.GetEnv(os.Getenv("HTTP_OTLP"), DefaultOtlpAddress),
			SwaggerAddress: config.GetEnv(os.Getenv("DRIVER_HTTP_SWAGGER_ADDRESS"), DefaultSwaggerAddress),
		},

		Mongo: MongoConfig{
			Database:      config.GetEnv(os.Getenv("DRIVER_DB"), DefaultDatabase),
			Uri:           config.GetEnv(os.Getenv("DRIVER_MONGO_URI"), DefaultMongoUri),
			MigrationsDir: config.GetEnv(os.Getenv("DRIVER_MONGO_MIGRATIONS_DIR"), DefaultMigrationsDir),
			Migrate:       config.GetEnvBoolean(os.Getenv("DRIVER_DB_MIGRATE")),
			Populate:      config.GetEnvBoolean(os.Getenv("DRIVER_DB_POPULATE")),
		},

		Kafka: KafkaConfig{
			Brokers:  config.GetBrokers(os.Getenv("DRIVER_KAFKA_BROKERS"), DefaultKafkaBrokers),
			GroupID:  config.GetEnv(os.Getenv("DRIVER_KAFKA_GROUP_ID"), DefaultKafkaGroupID),
			Topic:    config.GetEnv(os.Getenv("DRIVER_KAFKA_TOPIC"), DefaultKafkaTopic),
			MaxBytes: config.GetEnv(os.Getenv("DRIVER_KAFKA_MAX_BYTES"), DefaultKafkaMaxBytes),
		},

		Location: LocationService{
			URL: config.GetEnv(os.Getenv("LOCATION_URL"), DefaultLocationService),
		},
	}, nil
}
