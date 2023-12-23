package events

import (
	"github.com/shbov/hse-go_final/internal/driver/model/event_type"
	"github.com/shbov/hse-go_final/internal/driver/model/trip"
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"time"
)

type CreatedTripData struct {
	TripId  string                 `json:"trip_id"`
	OfferId string                 `json:"offer_id"`
	Price   trip.Price             `json:"price"`
	Status  trip_status.TripStatus `json:"status"`
	From    trip.Coordinates       `json:"from"`
	To      trip.Coordinates       `json:"to"`
}

type DefaultEvent struct {
	Id              string               `json:"id"`
	Source          string               `json:"source"`
	Type            event_type.EventType `json:"type"`
	DataContentType string               `json:"datacontenttype"`
	Time            time.Time            `json:"time"`
	Data            struct {
		TripId string `json:"trip_id"`
	} `json:"data"`
}

type CreatedTripEvent struct {
	Id              string          `json:"id"`
	Source          string          `json:"source"`
	Type            string          `json:"type"`
	DataContentType string          `json:"datacontenttype"`
	Time            time.Time       `json:"time"`
	Data            CreatedTripData `json:"data"`
}
