package message_queue

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type MessageQueue interface {
	GetReader(ctx context.Context) *kafka.Reader
	CancelTrip(ctx context.Context, tripId string, reason string) error
	AcceptTrip(ctx context.Context, driverId string, tripId string) error
	StartTrip(ctx context.Context, tripId string) error
	EndTrip(ctx context.Context, tripId string) error
}
