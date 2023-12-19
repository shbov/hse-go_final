package model

import "time"

type Location struct {
	Id       int    `json:"id"`
	DriverId string `json:"driver_id"`

	Lat float64 `json:"lat"` // float - ok?
	Lng float64 `json:"lng"`

	CreatedAt time.Time `json:"created_at"` // time.Time - ok? mb timestamp (other libs)?
}
