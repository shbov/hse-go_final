package requests

type GetDriversByLocationReqBody struct {
	Lat    float32 `json:"lat"`
	Lng    float32 `json:"lng"`
	Radius float32 `json:"radius"`
}

func (body *GetDriversByLocationReqBody) Validate() bool {
	if (body.Lat <= 90.0) && (body.Lat >= -90.0) && (body.Lng <= 180.0) && (body.Lng >= -180.0) {
		return true
	}
	return false
}
