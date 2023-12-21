package triprepo

import (
	"context"
	"errors"
	"github.com/shbov/hse-go_final/internal/driver/model"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"log"
	"time"
)

var _ repo.Trip = (*tripRepo)(nil)

type tripRepo struct {
	logger *zap.Logger
	db     *mongo.Database
}

const (
	timeoutQuery          = time.Second * 10
	collectionCreatedCode = 48
	tripCollection        = "trip"
)

func (r *tripRepo) GetTripsByUserId(ctx context.Context, userId string) ([]model.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	var trips []model.Trip

	coll := r.db.Collection(tripCollection)
	filter := bson.D{{"driver_id", userId}}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx) // ???

	if err = cursor.All(context.TODO(), &trips); err != nil {
		return nil, err
	}

	return trips, nil

}

func (r *tripRepo) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*model.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	var trip model.Trip

	coll := r.db.Collection(tripCollection)
	filter := bson.D{{"driver_id", userId}, {"id", tripId}}

	err := coll.FindOne(ctx, filter).Decode(trip)
	if err != nil {
		return nil, err
	}

	return &trip, nil
}

func New(db *mongo.Database, logger *zap.Logger) (repo.Trip, error) {
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

	log.Println("repo successfully created")
	return &tripRepo{
		logger: logger,
		db:     db,
	}, nil
}
