package service

import "context"

type Driver interface {
	SendTripInvitation(ctx context.Context, driverId string, TripId string) error
}
