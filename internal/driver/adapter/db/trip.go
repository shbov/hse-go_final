package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

const (
	timeoutQuery          = time.Second * 10
	collectionCreatedCode = 48
	tripCollection        = "trip"
)

type TripRepositoryImpl struct {
	logger *zap.Logger
	db     *mongo.Database
}

func CreateTripRepo(logger *zap.Logger, db *mongo.Database) (*TripRepositoryImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if db != nil {
		err := db.CreateCollection(ctx, tripCollection)
		if err != nil {
			var we mongo.CommandError
			if errors.As(err, &we) {
				if we.Code == collectionCreatedCode {
					err = nil
				}
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return &TripRepositoryImpl{
		logger: logger,
		db:     db,
	}, nil
}
