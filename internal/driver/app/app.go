package app

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/juju/zaputil/zapctx"
	"github.com/opentracing/opentracing-go"
	"github.com/shbov/hse-go_final/internal/driver/config"
	"github.com/shbov/hse-go_final/internal/driver/httpadapter"
	"github.com/shbov/hse-go_final/internal/driver/message_queue/drivermq"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/repo/triprepo"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"github.com/shbov/hse-go_final/internal/driver/service/kafkalistenersvc"
	"github.com/shbov/hse-go_final/internal/driver/service/kafkasvc"
	"github.com/shbov/hse-go_final/internal/driver/service/tripsvc"
	"github.com/shbov/hse-go_final/pkg/mongo_migration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var _ App = (*app)(nil)

type app struct {
	config *config.Config

	tracer       opentracing.Tracer
	messageQueue service.KafkaService
	tripService  service.Trip
	listener     service.Listener
	httpAdapter  httpadapter.Adapter
}

func (a *app) Serve(ctx context.Context) error {
	lg := zapctx.Logger(ctx)
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		a.listener.Run(ctx, a.config.Location.URL)
	}()

	go func() {
		lg.Info("server successfully started at " + a.config.HTTP.ServeAddress)
		if err := a.httpAdapter.Serve(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Error(err.Error())
		}
	}()

	<-done

	a.Shutdown(ctx)
	lg.Info("server successfully stopped")
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

	tripRepo, err := triprepo.New(ctx, db)
	if err != nil {
		return nil, err
	}

	tripService := tripsvc.New(ctx, tripRepo)

	mq, err := drivermq.New(&config.Kafka, lg)
	if err != nil {
		return nil, err
	}
	kafkaService := kafkasvc.New(ctx, mq)
	if err != nil {
		return nil, err
	}

	listener := kafkalistenersvc.New(ctx, tripService, kafkaService)

	a := &app{
		config:       config,
		tripService:  tripService,
		messageQueue: kafkaService,
		listener:     listener,
		httpAdapter:  httpadapter.New(ctx, &config.HTTP, kafkaService, tripService),
	}

	lg.Info("app successfully created")
	return a, nil
}

func initDB(ctx context.Context, config *config.Config) (*mongo.Database, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Mongo.Uri))
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
	if config.Mongo.Migrate {
		migrationSvc := mongo_migration.NewMigrationsService(lg, database)
		err = migrationSvc.RunMigrations(config.Mongo.MigrationsDir)
		if err != nil {
			return nil, err // fmt.Errorf("run migrations failed")
		}

		lg.Info("mongo db migrations finished")
	}

	if config.Mongo.Populate {
		result, err := database.Collection("trip").InsertMany(ctx, trip.FakeTrips)
		if err != nil {
			return nil, err
		}

		for _, id := range result.InsertedIDs {
			fmt.Printf("Inserted document with _id: %v\n", id)
		}
	}

	return database, nil
}
