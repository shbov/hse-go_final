package requests

type SetDriverLocationBody struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

func (body *SetDriverLocationBody) Validate() bool {
	if (body.Lat <= 90.0) && (body.Lat >= -90.0) && (body.Lng <= 180.0) && (body.Lng >= -180.0) {
		return true
	}
	return false
}
