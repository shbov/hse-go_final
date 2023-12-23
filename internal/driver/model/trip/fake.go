package trip

import (
	"github.com/shbov/hse-go_final/internal/driver/model/trip_status"
)

var tripExample1 = Trip{
	Id:       "17d52b7c-6f6f-462c-b5e7-98e7fdd14ca9",
	DriverId: "",

	From: Coordinates{10.0, 10.0},
	To:   Coordinates{12.0, 12.0},

	Price:  Price{99.5, "RUB"},
	Status: trip_status.DRIVERSEARCH,
}

var tripExample2 = Trip{
	Id:       "34ecca26-f17d-491f-9d5d-ef4db5c60876",
	DriverId: "23613a20-5787-42ef-ab81-3524a8e0c33f",

	From: Coordinates{44.0, 38.0},
	To:   Coordinates{45.0, 37.0},

	Price:  Price{1235.51, "RUB"},
	Status: trip_status.DRIVERFOUND,
}

var tripExample3 = Trip{
	Id:       "34ecca26-f07d-491f-9d5d-ef4db5c60876",
	DriverId: "23613a20-5787-42ef-ab11-3524a8e0c33f",

	From: Coordinates{44.0, 38.0},
	To:   Coordinates{45.0, 37.0},

	Price:  Price{1235.51, "RUB"},
	Status: trip_status.ENDED,
}

var FakeTrips = []interface{}{tripExample1, tripExample2, tripExample3}
