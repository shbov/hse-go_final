package trip

import (
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
)

type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type Price struct {
	Amount   float64 `json:"amount" bson:"amount"`
	Currency string  `json:"currency" bson:"currency"`
}

type Trip struct {
	Id       string `json:"id" bson:"id"`
	DriverId string `json:"driver_id" bson:"driver_id"`

	From Coordinates `json:"from" bson:"from"`
	To   Coordinates `json:"to" bson:"to"`

	Price Price `json:"price" bson:"price"`

	Status trip_status.TripStatus `json:"status" bson:"status"`
}
