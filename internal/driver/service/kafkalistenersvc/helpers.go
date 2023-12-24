package kafkalistenersvc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/model/events"
	"github.com/shbov/hse-go_final/internal/driver/model/requests"
	"github.com/shbov/hse-go_final/internal/driver/model/responses"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/service/driversvc"
	"net/http"
)

func ParseEventCreate(m kafka.Message) (*trip.Trip, error) {
	var createEvent events.CreatedTripEvent
	if err := json.Unmarshal(m.Value, &createEvent); err != nil {
		return nil, err
	}
	tripToSave := trip.Trip{
		Id:       createEvent.Data.TripId,
		DriverId: "",
		From:     createEvent.Data.From,
		To:       createEvent.Data.To,
		Price:    createEvent.Data.Price,
		Status:   createEvent.Data.Status,
	}
	return &tripToSave, nil
}

func SendTripInvitationsToDrivers(ctx context.Context, lat float64, lng float64, radius float64, locationURL string, tripId string) error {
	body := requests.GetDriversBody{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
	}

	req, err := http.NewRequest(http.MethodGet, locationURL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("lat", fmt.Sprintf("%f", body.Lat))
	q.Add("lng", fmt.Sprintf("%f", body.Lng))
	q.Add("radius", fmt.Sprintf("%f", body.Radius))
	req.URL.RawQuery = q.Encode()

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	var drivers responses.DriverInfoArray
	if err := json.NewDecoder(r.Body).Decode(&drivers); err != nil {
		return err
	}

	for _, driver := range drivers {
		if err := driversvc.New().SendTripInvitation(ctx, driver.DriverId, tripId); err != nil {
			return err
		}
	}
	return nil
}
