package responses

type DriverInfoArray []struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	DriverId string  `json:"driver_id"`
	ID       int64   `json:"id"`
}
