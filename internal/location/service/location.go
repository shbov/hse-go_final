package service

import (
	"github.com/shbov/hse-go_final/internal/location/model"
	"golang.org/x/net/context"
)

type Location interface {
	GetDriversInLocation(ctx context.Context, centerLat float32, centerLng float32, radius float32) ([]model.Location, error)
	SetLocationByDriverId(ctx context.Context, driverId string, lat float32, lng float32) error
}
