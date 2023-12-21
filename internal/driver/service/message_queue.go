package service

import (
	"golang.org/x/net/context"
)

type MessageQueue interface {
	CancelTrip(ctx context.Context, userId string, tripId string) error
	AcceptTrip(ctx context.Context, userId string, tripId string) error
	StartTrip(ctx context.Context, userId string, tripId string) error
	EndTrip(ctx context.Context, userId string, tripId string) error
}
