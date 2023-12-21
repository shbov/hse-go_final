package kafkasvc

import (
	"context"
	"github.com/shbov/hse-go_final/internal/driver/message_queue"
	"github.com/shbov/hse-go_final/internal/driver/service"
	"log"
)

var _ service.KafkaService = (*driverService)(nil)

type driverService struct {
	mq message_queue.MessageQueue
}

func (ls *driverService) CancelTrip(ctx context.Context, tripId string, reason string) error {
	err := ls.mq.CancelTrip(ctx, tripId, reason)
	return err
}
func (ls *driverService) AcceptTrip(ctx context.Context, driverId string, tripId string) error {
	err := ls.mq.AcceptTrip(ctx, driverId, tripId)
	return err
}
func (ls *driverService) StartTrip(ctx context.Context, tripId string) error {
	err := ls.mq.StartTrip(ctx, tripId)
	return err
}
func (ls *driverService) EndTrip(ctx context.Context, tripId string) error {
	err := ls.mq.EndTrip(ctx, tripId)
	return err
}

func New(mq message_queue.MessageQueue) service.KafkaService {
	s := &driverService{
		mq: mq,
	}

	log.Println("kafka service successfully created")
	return s
}
