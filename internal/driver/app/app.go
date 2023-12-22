package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/opentracing/opentracing-go"
	"github.com/shbov/hse-go_final/internal/driver/config"

	"github.com/shbov/hse-go_final/internal/driver/httpadapter"
	"github.com/shbov/hse-go_final/internal/driver/message_queue/drivermq"
	"github.com/shbov/hse-go_final/internal/driver/repo/triprepo"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"github.com/shbov/hse-go_final/internal/driver/service/tripsvc"
	"github.com/shbov/hse-go_final/pkg/mongo_migration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var _ App = (*app)(nil)

type app struct {
	config *config.Config

	tracer       opentracing.Tracer
	messageQueue service.KafkaService
	tripService  service.Trip
	httpAdapter  httpadapter.Adapter
}

func (a *app) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("server successfully started at " + a.config.HTTP.ServeAddress)
		if err := a.httpAdapter.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Fatal(err.Error())
		}
	}()

	<-done

	a.Shutdown(ctx)
	log.Println("server successfully stopped")
	return nil
}

func (a *app) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, a.config.App.ShutdownTimeout)
	defer cancel()

	a.httpAdapter.Shutdown(ctx)
}

func New(ctx context.Context, config *config.Config) (App, error) {
	lg := zapctx.Logger(ctx)

	db, err := initDB(ctx, config)
	if err != nil {
		return nil, err
	}

	tripRepo, err := triprepo.New(db)
	if err != nil {
		return nil, err
	}

	tripService := tripsvc.New(tripRepo)
	messageQueue, err := drivermq.New(&config.Kafka, lg)
	if err != nil {
		return nil, err
	}

	a := &app{
		config:       config,
		tripService:  tripService,
		messageQueue: messageQueue,
		httpAdapter:  httpadapter.New(&config.HTTP, messageQueue, tripService),
	}

	log.Println("app successfully created")
	return a, nil
}

func initDB(ctx context.Context, config *config.Config) (*mongo.Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.Uri))
	if err != nil {
		return nil, fmt.Errorf("new mongo client create error: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("new mongo primary node connect error: %w", err)
	}

	database := client.Database(config.Mongo.Database)
	lg := zapctx.Logger(ctx)

	// uncomment if migration is needed
	migrationSvc := mongo_migration.NewMigrationsService(lg, database)
	err = migrationSvc.RunMigrations(config.Mongo.MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("run migrations failed")
	}
	log.Println("mongo db migrations finished")

	return database, nil
}
