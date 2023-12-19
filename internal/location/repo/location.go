package repo

import (
	"context"

	"github.com/shbov/hse-go_final/internal/location/model"
)

type Location interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddLocation(ctx context.Context, driverId string, lat float64, lng float64) error
	GetLocation(ctx context.Context, centerLat float64, centerLng float64, radius float64) (*model.Location, error)
	GetLocationByDriverId(ctx context.Context, driverId string) (*model.Location, error)
}
