package responses

type DriverInfo struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	DriverId string  `json:"id"`
}
