package repo

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
)

type Trip interface {
	AddTrip(ctx context.Context, trip trip.Trip) error
	GetTripsByUserId(ctx context.Context, userId string) ([]trip.Trip, error)
	GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*trip.Trip, error)
}
