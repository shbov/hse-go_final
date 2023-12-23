package kafkalistenersvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/model/events"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.Listener = (*kafkaListener)(nil)

const (
	eventTypeCreated string = "trip.event.created"
)

type kafkaListener struct {
	tripService  service.Trip
	kafkaService service.KafkaService
}

func (kl *kafkaListener) Run(ctx context.Context) {
	lg := zapctx.Logger(ctx)
	reader := kl.kafkaService.GetReader(ctx)
	for {
		select {
		case <-ctx.Done(): // will execute if cancel func is called.
			if err := reader.Close(); err != nil {
				lg.Fatal(fmt.Sprintf("failed to close reader: %s\n", err))
			}
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				break
			}
			var event events.DefaultEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				lg.Fatal(fmt.Sprintf("failed to unmarshal event: %s\n", err))
				return
			}
			if event.Type == eventTypeCreated {
				var createEvent events.CreatedTripEvent
				if err := json.Unmarshal(m.Value, &createEvent); err != nil {
					lg.Fatal(fmt.Sprintf("failed to unmarshal event: %s\n", err))
					return
				}
				tripToSave := trip.Trip{
					Id:       createEvent.Data.TripId,
					DriverId: "",
					From:     trip.Coordinates(createEvent.Data.From),
					To:       trip.Coordinates(createEvent.Data.To),
					Price:    trip.Price(createEvent.Data.Price),
					Status:   createEvent.Data.Status,
				}
				if err := kl.tripService.AddTrip(ctx, tripToSave); err != nil {
					lg.Fatal(fmt.Sprintf("failed to save trip: %s\n", err))
				}
			} else {
				var status trip_status.TripStatus
				switch event.Type {
				case "trip.event.accepted":
					status = trip_status.ACCEPTED
				case "trip.event.cancelled":
					status = trip_status.CANCELED
				case "trip.event.ended":
					status = trip_status.ENDED
				case "trip.event.started":
					status = trip_status.STARTED
				}
				if err := kl.tripService.ChangeTripStatus(ctx, event.Data.TripId, status); err != nil {
					lg.Fatal(fmt.Sprintf("failed to update trip status: %s\n", err))
				}
			}
		}
	}
}

func New(ctx context.Context, ts service.Trip, ks service.KafkaService) service.Listener {
	lg := zapctx.Logger(ctx)
	kl := &kafkaListener{
		tripService:  ts,
		kafkaService: ks,
	}

	lg.Info("kafka listener successfully created")

	return kl
}
