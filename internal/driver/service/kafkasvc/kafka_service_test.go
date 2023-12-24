package kafkasvc

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/config"
	"github.com/shbov/hse-go_final/internal/driver/message_queue/kafkamock"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReader(t *testing.T) {
	mockKafka := new(kafkamock.KafkaMock)

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		log.Fatal(err.Error())
	}

	rc := kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		GroupID: cfg.Kafka.GroupID,
		Topic:   cfg.Kafka.Topic,
	}

	ctx := context.Background()

	mockKafka.On("GetReader", ctx).Return(kafka.NewReader(rc))
	service := New(ctx, mockKafka)

	reader := service.GetReader(ctx)
	assert.NotNil(t, reader, "Expected result not to be nil")

	mockKafka.AssertExpectations(t)
}

func TestCancelTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	service := New(context.Background(), mockKafka)
	tripId := "trip-1"
	reason := "no fuel"

	ctx := context.Background()

	mockKafka.On("CancelTrip", ctx, tripId, reason).Return(nil)

	// Test
	err := service.CancelTrip(ctx, tripId, reason)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}

func TestAcceptTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	driverId := "driver-1"
	tripId := "trip-1"

	service := New(context.Background(), mockKafka)

	ctx := context.Background()

	mockKafka.On("AcceptTrip", ctx, driverId, tripId).Return(nil)

	// Test
	err := service.AcceptTrip(ctx, driverId, tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}

func TestStartTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	tripId := "trip-1"

	service := New(context.Background(), mockKafka)

	ctx := context.Background()

	mockKafka.On("StartTrip", ctx, tripId).Return(nil)

	// Test
	err := service.StartTrip(ctx, tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}

func TestEndTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	tripId := "trip-1"

	service := New(context.Background(), mockKafka)

	ctx := context.Background()

	mockKafka.On("EndTrip", ctx, tripId).Return(nil)

	// Test
	err := service.EndTrip(ctx, tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}
