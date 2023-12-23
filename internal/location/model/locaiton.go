package model

type Location struct {
	Id       int    `json:"id"`
	DriverId string `json:"driver_id"`

	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
