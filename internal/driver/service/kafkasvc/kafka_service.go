package kafkasvc

import (
	"context"
	"github.com/juju/zaputil/zapctx"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/message_queue"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.KafkaService = (*driverService)(nil)

type driverService struct {
	mq message_queue.MessageQueue
}

func (ds *driverService) GetReader(ctx context.Context) *kafka.Reader {
	return ds.mq.GetReader(ctx)
}

func (ds *driverService) CancelTrip(ctx context.Context, tripId string, reason string) error {
	err := ds.mq.CancelTrip(ctx, tripId, reason)
	return err
}

func (ds *driverService) AcceptTrip(ctx context.Context, driverId string, tripId string) error {
	err := ds.mq.AcceptTrip(ctx, driverId, tripId)
	return err
}

func (ds *driverService) StartTrip(ctx context.Context, tripId string) error {
	err := ds.mq.StartTrip(ctx, tripId)
	return err
}

func (ds *driverService) EndTrip(ctx context.Context, tripId string) error {
	err := ds.mq.EndTrip(ctx, tripId)
	return err
}

func New(ctx context.Context, mq message_queue.MessageQueue) service.KafkaService {
	lg := zapctx.Logger(ctx)

	s := &driverService{
		mq: mq,
	}

	lg.Info("kafka service successfully created")
	return s
}
