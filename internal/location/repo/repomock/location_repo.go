package repomock

import (
	"context"

	"github.com/shbov/hse-go_final/internal/location/model"
	"github.com/shbov/hse-go_final/internal/location/repo"

	"github.com/stretchr/testify/mock"
)

var _ repo.Location = (*LocationMock)(nil)

type LocationMock struct {
	mock.Mock
}

func (l *LocationMock) AddLocation(ctx context.Context, driverId string, lat float32, lng float32) error {
	args := l.Called(ctx, driverId, lat, lng)
	return args.Error(0)
}

func (l *LocationMock) GetDriversInLocation(ctx context.Context, centerLat float32, centerLng float32, radius float32) ([]model.Location, error) {
	args := l.Called(ctx, centerLat, centerLng, radius)
	return args.Get(0).([]model.Location), args.Error(1)
}

func (l *LocationMock) SetLocationByDriverId(ctx context.Context, driverId string, lat float32, lng float32) error {
	args := l.Called(ctx, driverId, lat, lng)
	return args.Error(0)
}
