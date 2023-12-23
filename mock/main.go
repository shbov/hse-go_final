package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/model/events"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"log"
	"time"
)

var async = flag.Bool("a", false, "use async")

func writeToKafka(ctx context.Context, writer *kafka.Writer, event any) error {
	m, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if err := writer.WriteMessages(ctx, kafka.Message{Value: m}); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	logger := log.Default()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:     []string{"localhost:29092"},
		Topic:       "driver",
		Async:       *async,
		Logger:      kafka.LoggerFunc(logger.Printf),
		ErrorLogger: kafka.LoggerFunc(logger.Printf),
	})
	defer writer.Close()

	eventCreate1 := events.CreatedTripEvent{
		Id:              "e322a084-ca0f-4f4a-91cf-cae27946954b",
		Source:          "/trip",
		Type:            "event.trip.created",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: events.CreatedTripData{
			TripId:  "06b695a3-5c9d-4d40-aff8-ba7943145370",
			OfferId: "de87cd91-f109-48f0-afc1-64cdc7579fe3",
			Price: trip.Price{
				Amount:   524.1234,
				Currency: "KAFKA_CUR",
			},
			Status: "CREATED_FROM_KAFKA",
			From: trip.Coordinates{
				Lat: 43.23,
				Lng: 54.22,
			},
			To: trip.Coordinates{
				Lat: 65.12,
				Lng: 73.73,
			},
		},
	}

	if err := writeToKafka(ctx, writer, eventCreate1); err != nil {
		log.Fatal("err write to kafka:", err)
	}

	tripToAddFromKafka2 := events.CreatedTripEvent{
		Id:              "0a580514-c9df-4ffe-a24a-85c21c84d7d6",
		Source:          "/trip",
		Type:            "event.trip.created",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: events.CreatedTripData{
			TripId:  "37b225c8-c0b5-4848-90ea-3a5860cff532",
			OfferId: "bf79b2a9-ea41-45b0-bf13-2417757e38c8",
			Price: trip.Price{
				Amount:   6912.185,
				Currency: "KAFKA_CUR",
			},
			Status: "CREATED",
			From: trip.Coordinates{
				Lat: 83.58,
				Lng: 138.12,
			},
			To: trip.Coordinates{
				Lat: 72.12,
				Lng: -12.73,
			},
		},
	}

	if err := writeToKafka(ctx, writer, tripToAddFromKafka2); err != nil {
		log.Fatal("err write to kafka:", err)
	}

	changeStatus := events.DefaultEvent{
		Id:              "e50611ac-af1f-4a47-a38d-ad9f6dfa547d",
		Source:          "/trip",
		Type:            "trip.event.started",
		DataContentType: "application/json",
		Time:            time.Time{},
		Data: struct {
			TripId string `json:"trip_id"`
		}{TripId: "37b225c8-c0b5-4848-90ea-3a5860cff532"},
	}

	if err := writeToKafka(ctx, writer, changeStatus); err != nil {
		log.Fatal("err write to kafka:", err)
	}

}
