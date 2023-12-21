package service

import (
	"golang.org/x/net/context"
)

type MessageQueue interface {
	CancelTrip(ctx context.Context, tripId string, reason string) error
	AcceptTrip(ctx context.Context, driverId string, tripId string) error
	StartTrip(ctx context.Context, tripId string) error
	EndTrip(ctx context.Context, tripId string) error
}
