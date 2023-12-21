package service

import (
	"github.com/shbov/hse-go_final/internal/driver/model"
	"golang.org/x/net/context"
)

type Driver interface {
	GetTrips(ctx context.Context, userId string) ([]model.Location, error)
	SetLocationByDriverId(ctx context.Context, driverId string, lat float32, lng float32) error
}
