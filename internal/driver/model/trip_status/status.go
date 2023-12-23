package trip_status

type TripStatus string

const (
	CANCELED     TripStatus = "CANCELED"
	DRIVERFOUND  TripStatus = "DRIVER_FOUND"
	DRIVERSEARCH TripStatus = "DRIVER_SEARCH"
	ENDED        TripStatus = "ENDED"
	ONPOSITION   TripStatus = "ON_POSITION"
	STARTED      TripStatus = "STARTED"
)
