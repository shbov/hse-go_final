package model

type Location struct {
	Id       int    `json:"id"`
	DriverId string `json:"driver_id"`

	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}
