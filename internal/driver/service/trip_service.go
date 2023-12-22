package service

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/model"
)

type Trip interface {
	AddTrip(ctx context.Context, trip model.Trip) error
	GetTripsByUserId(ctx context.Context, userId string) ([]model.Trip, error)
	GetTripByUserIdTripId(ctx context.Context, userId string, tripId string) (*model.Trip, error)
}
