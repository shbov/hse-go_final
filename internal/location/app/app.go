package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/juju/zaputil/zapctx"
	"github.com/opentracing/opentracing-go"
	"github.com/shbov/hse-go_final/internal/location/httpadapter"
	"github.com/shbov/hse-go_final/internal/location/repo/locationrepo"
	"github.com/shbov/hse-go_final/internal/location/service"
	"github.com/shbov/hse-go_final/internal/location/service/locationsvc"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type app struct {
	config *Config

	tracer          opentracing.Tracer
	locationService service.Location
	httpAdapter     httpadapter.Adapter
}

func (a *app) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := a.httpAdapter.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown()

	return nil
}

func (a *app) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(ctx context.Context, config *Config) (App, error) {
	pgxPool, err := initDB(context.Background(), &config.Database)
	if err != nil {
		return nil, err
	}

	locationRepo, err := locationrepo.New(pgxPool)
	if err != nil {
		return nil, err
	}

	locationService := locationsvc.New(locationRepo)

	a := &app{
		config:          config,
		locationService: locationService,
		httpAdapter:     httpadapter.New(&config.HTTP, httpadapter.OurService{}),
	}

	return a, nil
}

func initDB(ctx context.Context, config *DatabaseConfig) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// migrations
	m, err := migrate.New(config.MigrationsDir, config.DSN)
	if err != nil {
		return nil, err
	}

	//if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
	//	return nil, err
	//}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return pool, nil
}
