package events

import (
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
	"time"
)

type DefaultEvent struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            struct {
		TripId string `json:"trip_id"`
	} `json:"data"`
}

type CreatedTripEvent struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            struct {
		TripId  string `json:"trip_id"`
		OfferId string `json:"offer_id"`
		Price   struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		} `json:"price"`
		Status trip_status.TripStatus `json:"status"`
		From   struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"from"`
		To struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"to"`
	} `json:"data"`
}
