package repo

import (
	"context"

	"github.com/shbov/hse-go_final/internal/location/model"
)

type Location interface {
	GetDriversInLocation(ctx context.Context, centerLat float64, centerLng float64, radius float64) (*model.Location, error)
	SetLocationByDriverId(ctx context.Context, driverId string, lat float64, lng float64) error
}
