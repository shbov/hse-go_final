package driversvc

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.Driver = (*driverService)(nil)

type driverService struct{}

func (ds *driverService) SendTripInvitation(ctx context.Context, driverId string, TripId string) error {
	return nil
}

func New() service.Driver {
	return &driverService{}
}
