package model

import "time"

type Location struct {
	Id       int    `json:"id"`
	DriverId string `json:"driver_id"`

	Lat float32 `json:"lat"` // float - ok?
	Lng float32 `json:"lng"`

	CreatedAt time.Time `json:"created_at"` // time.Time - ok? mb timestamp (other libs)?
}
