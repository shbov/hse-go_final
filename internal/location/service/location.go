package service

import (
	"golang.org/x/net/context"
)

type Location interface {
	SetLocationByDriverId(ctx context.Context, driverId string, lat float64, lng float64) error
}
