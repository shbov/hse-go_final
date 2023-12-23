package drivermq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/config"
	"github.com/shbov/hse-go_final/internal/driver/message_queue"
	"github.com/shbov/hse-go_final/internal/driver/model/commands"
	"go.uber.org/zap"
	"time"
)

const (
	source          = "/driver"
	dataContentType = "application/json"
	typePrefix      = "trip.command."
)

var _ message_queue.MessageQueue = (*driverKafka)(nil)

type driverKafka struct {
	wc kafka.WriterConfig
	rc kafka.ReaderConfig
}

type logWrap struct {
	l *zap.Logger
}

func (wrap logWrap) logf(msg string, a ...interface{}) {
	wrap.l.Info(fmt.Sprintf(msg, a))
}

func (d *driverKafka) GetReader(ctx context.Context) *kafka.Reader {
	return kafka.NewReader(d.rc)
}

func (d *driverKafka) CancelTrip(ctx context.Context, tripId string, reason string) error {
	writer := kafka.NewWriter(d.wc)
	id := uuid.New().String()
	msg := commands.CancelMessage{
		Id:              id,
		Source:          source,
		Type:            typePrefix + "cancel",
		DataContentType: dataContentType,
		Time:            time.Now(),
		Data: struct {
			TripId string `json:"trip_id"`
			Reason string `json:"reason"`
		}{TripId: tripId, Reason: reason},
	}

	parsedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{Value: parsedMsg})
	if err != nil {
		return err
	}

	lg := zapctx.Logger(ctx)
	if err := writer.Close(); err != nil {
		lg.Error(fmt.Sprintf("failed to close writer: %s\n", err))
	}

	return nil
}

func (d *driverKafka) AcceptTrip(ctx context.Context, driverId string, tripId string) error {
	writer := kafka.NewWriter(d.wc)
	id := uuid.New().String()
	msg := commands.AcceptMessage{
		Id:              id,
		Source:          source,
		Type:            typePrefix + "accept",
		DataContentType: dataContentType,
		Time:            time.Now(),
		Data: struct {
			TripId   string `json:"trip_id"`
			DriverId string `json:"driver_id"`
		}{TripId: tripId, DriverId: driverId},
	}

	parsedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{Value: parsedMsg})
	if err != nil {
		return err
	}

	lg := zapctx.Logger(ctx)
	if err := writer.Close(); err != nil {
		lg.Error(fmt.Sprintf("failed to close writer: %s\n", err))
	}

	return nil
}
func (d *driverKafka) StartTrip(ctx context.Context, tripId string) error {
	writer := kafka.NewWriter(d.wc)
	id := uuid.New().String()
	msg := commands.StartEndMessage{
		Id:              id,
		Source:          source,
		Type:            typePrefix + "start",
		DataContentType: dataContentType,
		Time:            time.Now(),
		Data: struct {
			TripId string `json:"trip_id"`
		}{TripId: tripId},
	}

	parsedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{Value: parsedMsg})
	if err != nil {
		return err
	}

	lg := zapctx.Logger(ctx)
	if err := writer.Close(); err != nil {
		lg.Error(fmt.Sprintf("failed to close writer: %s\n", err))
	}

	return nil
}
func (d *driverKafka) EndTrip(ctx context.Context, tripId string) error {
	writer := kafka.NewWriter(d.wc)
	id := uuid.New().String()
	msg := commands.StartEndMessage{
		Id:              id,
		Source:          source,
		Type:            typePrefix + "end",
		DataContentType: dataContentType,
		Time:            time.Now(),
		Data: struct {
			TripId string `json:"trip_id"`
		}{TripId: tripId},
	}

	parsedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{Value: parsedMsg})
	if err != nil {
		return err
	}

	lg := zapctx.Logger(ctx)
	if err := writer.Close(); err != nil {
		lg.Error(fmt.Sprintf("failed to close writer: %s\n", err))
	}
	return nil
}

func New(conf *config.KafkaConfig, lg *zap.Logger) (message_queue.MessageQueue, error) {
	w := kafka.WriterConfig{
		Brokers:  conf.Brokers,
		Topic:    conf.Topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   kafka.LoggerFunc(logWrap{lg}.logf),
	}

	r := kafka.ReaderConfig{
		Brokers: conf.Brokers,
		GroupID: conf.GroupID,
		Topic:   conf.Topic,
	}

	d := &driverKafka{
		wc: w,
		rc: r,
	}

	lg.Info(fmt.Sprintf("message_queue successfully created; listening addresses: %s", conf.Brokers))
	return d, nil
}
