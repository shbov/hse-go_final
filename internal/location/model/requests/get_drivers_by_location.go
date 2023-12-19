package requests

type GetDriversByLocationReqBody struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
}

func (body *GetDriversByLocationReqBody) Validate() bool {
	if (body.Lat <= 90.0) && (body.Lat >= -90.0) && (body.Lng <= 180.0) && (body.Lng >= -180.0) {
		return true
	}
	return false
}
