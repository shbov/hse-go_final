package kafkalistenersvc

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/shbov/hse-go_final/internal/driver/model/events"
	"github.com/shbov/hse-go_final/internal/driver/model/requests"
	"github.com/shbov/hse-go_final/internal/driver/model/responses"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/service/driversvc"
	"io/ioutil"
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

	bodyJSONed, err := json.Marshal(body)
	bodyEncoded := bytes.NewReader(bodyJSONed)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, locationURL, bodyEncoded)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var drivers []responses.DriverInfo
	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(reqBody, &drivers); err != nil {
		return err
	}

	for _, driver := range drivers {
		if err := driversvc.New().SendTripInvitation(ctx, driver.DriverId, tripId); err != nil {
			return err
		}
	}
	return nil
}
