package triprepo

import (
	"context"
	"errors"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var _ repo.Trip = (*tripRepo)(nil)

type tripRepo struct {
	db *mongo.Database
}

const (
	timeoutQuery          = time.Second * 10
	collectionCreatedCode = 48
	tripCollection        = "trip"
)

func (r *tripRepo) AddTrip(ctx context.Context, trip trip.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	coll := r.db.Collection(tripCollection)

	_, err := coll.InsertOne(ctx, trip)
	if err != nil {
		return err
	}

	return nil
}

func (r *tripRepo) GetTripsByUserId(ctx context.Context, userId string) ([]trip.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	var trips []trip.Trip

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

func (r *tripRepo) GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*trip.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutQuery)
	defer cancel()

	var trip trip.Trip

	coll := r.db.Collection(tripCollection)
	filter := bson.D{{"driver_id", userId}, {"id", tripId}}

	err := coll.FindOne(ctx, filter).Decode(&trip)
	if err != nil {
		return nil, err
	}

	return &trip, nil
}

func New(ctx context.Context, db *mongo.Database) (repo.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	lg := zapctx.Logger(ctx)
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

	lg.Info("repo successfully created")
	return &tripRepo{
		db: db,
	}, nil
}
