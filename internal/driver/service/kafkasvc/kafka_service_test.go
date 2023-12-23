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

	mockKafka.On("AddTrip", context.Background()).Return(kafka.NewReader(rc))
	kafkaService := New(context.Background(), mockKafka)

	reader := kafkaService.GetReader(context.Background())
	assert.NotNil(t, reader, "Expected result not to be nil")
	assert.Equal(t, reader, kafka.NewReader(rc), "Expected reader with the same config")

	mockKafka.AssertExpectations(t)
}

func TestCancelTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	service := New(context.Background(), mockKafka)
	tripId := "trip-1"
	reason := "no fuel"

	mockKafka.On("CancelTrip", context.Background(), tripId, reason).Return(nil)

	// Test
	err := service.CancelTrip(context.Background(), tripId, reason)

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
	mockKafka.On("CancelTrip", context.Background(), tripId, driverId).Return(nil)

	// Test
	err := service.AcceptTrip(context.Background(), driverId, tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}

func TestStartTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	driverId := "driver-1"
	tripId := "trip-1"

	service := New(context.Background(), mockKafka)
	mockKafka.On("CancelTrip", context.Background(), tripId, driverId).Return(nil)

	// Test
	err := service.StartTrip(context.Background(), tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}

func TestEndTrip(t *testing.T) {
	// Setup
	mockKafka := new(kafkamock.KafkaMock)

	driverId := "driver-1"
	tripId := "trip-1"

	service := New(context.Background(), mockKafka)
	mockKafka.On("CancelTrip", context.Background(), tripId, driverId).Return(nil)

	// Test
	err := service.EndTrip(context.Background(), tripId)

	// Assertions
	assert.NoError(t, err, "Expected no error")

	mockKafka.AssertExpectations(t)
}
