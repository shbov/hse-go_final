package kafkalistenersvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	"github.com/shbov/hse-go_final/internal/driver/model/event_type"
	"github.com/shbov/hse-go_final/internal/driver/model/events"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"github.com/shbov/hse-go_final/internal/driver/service"
)

var _ service.Listener = (*kafkaListener)(nil)

type kafkaListener struct {
	tripService  service.Trip
	kafkaService service.KafkaService
}

const defaultSearchRadius float64 = 0.02

func (kl *kafkaListener) Run(ctx context.Context, locationURL string) {
	lg := zapctx.Logger(ctx)
	reader := kl.kafkaService.GetReader(ctx)

	for {
		select {
		case <-ctx.Done(): // will execute if cancel func is called.
			if err := reader.Close(); err != nil {
				lg.Error(fmt.Sprintf("failed to close reader: %s\n", err))
			}
			return
		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				lg.Error(fmt.Sprintf("failed to read event: %s\n", err))
				return
			}
			lg.Info("read new message from kafka")

			var event events.DefaultEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				lg.Error(fmt.Sprintf("failed to unmarshal event: %s\n", err))
				return
			}

			if event.Type == event_type.CREATED {
				tripToSave, err := ParseEventCreate(m)
				if err != nil {
					lg.Error(fmt.Sprintf("failed to parse event: %s\n", err))
				}

				if err := kl.tripService.AddTrip(ctx, *tripToSave); err != nil {

					lg.Error(fmt.Sprintf("failed to save trip: %s\n", err))
				}

				err = SendTripInvitationsToDrivers(
					ctx,
					tripToSave.From.Lat,
					tripToSave.From.Lng,
					defaultSearchRadius,
					locationURL,
					tripToSave.Id,
				)
				if err != nil {
					lg.Error(fmt.Sprintf("failed to send trip invitations: %s\n", err))
				}

			} else {
				var status trip_status.TripStatus
				switch event.Type {
				case event_type.ACCEPTED:
					status = trip_status.ACCEPTED
				case event_type.CANCELED:
					status = trip_status.CANCELED
				case event_type.ENDED:
					status = trip_status.ENDED
				case event_type.STARTED:
					status = trip_status.STARTED
				}

				if err := kl.tripService.ChangeTripStatus(ctx, event.Data.TripId, status); err != nil {
					lg.Error(fmt.Sprintf("failed to update trip status: %s\n", err))
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
