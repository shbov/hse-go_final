package service

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
)

type Trip interface {
	ChangeTripStatus(ctx context.Context, tripId string, status trip_status.TripStatus) error
	AddTrip(ctx context.Context, trip trip.Trip) error
	GetTripsByUserId(ctx context.Context, userId string) ([]trip.Trip, error)
	GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*trip.Trip, error)
}
