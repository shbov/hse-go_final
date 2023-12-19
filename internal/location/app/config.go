package app

import (
	"github.com/shbov/hse-go_final/internal/location/httpadapter"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	AppName                = "location"
	DefaultServeAddress    = "localhost:8080"
	DefaultShutdownTimeout = 20 * time.Second
	DefaultBasePath        = "/location/v1"
	DefaultDSN             = "dsn://"
	DefaultMigrationsDir   = "file://migrations/location"

	DefaultAccessTokenCookie  = "access_token"
	DefaultRefreshTokenCookie = "refresh_token"

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

func NewConfig(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	cnf := Config{
		App: AppConfig{
			ShutdownTimeout: DefaultShutdownTimeout,
			DSN:             DefaultDSN,
		},
		HTTP: httpadapter.Config{
			ServeAddress:       DefaultServeAddress,
			BasePath:           DefaultBasePath,
			AccessTokenCookie:  DefaultAccessTokenCookie,
			RefreshTokenCookie: DefaultRefreshTokenCookie,
			OtlpAddress:        DefaultOtlpAddress,
		},
		Database: DatabaseConfig{
			DSN:           DefaultDSN,
			MigrationsDir: DefaultMigrationsDir,
		},
	}

	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}
