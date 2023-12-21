package app

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	config *Config

	tracer opentracing.Tracer

	mongoClient *mongo.Client

	log *zap.Logger

	//httpAdapter httpadapter.Adapter
}

func (a *app) Serve(ctx context.Context) error {
	//lg := zapctx.Logger(ctx)
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.log.Info("Connecting to mongo...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(a.config.Mongo.Uri))
	if err != nil {
		return fmt.Errorf("new mongo client create error: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("new mongo primary node connect error: %w", err)
	}

	//a.mongoClient = client
	//database := client.Database(a.config.Mongo.Database)

	//if a.config.Migrations.Enabled {
	//	migrationSvc := migration.NewMigrationsService(a.log, database)
	//	err = migrationSvc.RunMigrations(a.config.Migrations.Path)
	//	if err != nil {
	//		return fmt.Errorf("run migration failed")
	//	}
	//}
	//
	//tripRepo, err := db.CreateTripRepo(a.log, database)
	//if err != nil {
	//	return fmt.Errorf("url repo create failed: %w", err)
	//}
	//go func() {
	//	if err := a.httpAdapter.Serve(ctx); err != nil && err != http.ErrServerClosed {
	//		lg.Fatal(err.Error())
	//	}
	//}()

	<-done

	a.Shutdown()

	return nil
}

func (a *app) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), a.config.App.ShutdownTimeout)
	defer cancel()

	//a.httpAdapter.Shutdown(ctx)
}

func New(ctx context.Context, config *Config) (App, error) {
	a := &app{
		config: config,
		//httpAdapter: nil,
	}
	return a, nil
}
