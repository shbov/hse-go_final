package kafkamock

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/message_queue"
	"github.com/stretchr/testify/mock"
)

var _ message_queue.MessageQueue = (*KafkaMock)(nil)

type KafkaMock struct {
	mock.Mock
}

func (r *KafkaMock) EndTrip(ctx context.Context, tripId string) error {
	args := r.Called(ctx, tripId)
	return args.Error(0)
}

func (r *KafkaMock) StartTrip(ctx context.Context, tripId string) error {
	args := r.Called(ctx, tripId)
	return args.Error(0)
}

func (r *KafkaMock) AcceptTrip(ctx context.Context, driverId string, tripId string) error {
	args := r.Called(ctx, driverId, tripId)
	return args.Error(0)
}

func (r *KafkaMock) CancelTrip(ctx context.Context, tripId string, reason string) error {
	args := r.Called(ctx, tripId, reason)
	return args.Error(0)
}

func (r *KafkaMock) GetReader(ctx context.Context) *kafka.Reader {
	args := r.Called(ctx)
	return args.Get(0).(*kafka.Reader)
}
