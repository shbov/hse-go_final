package model

import "github.com/shopspring/decimal"

type coordinates struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type price struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}

type Trip struct {
	Id       string `json:"id"`
	DriverId string `json:"driver_id"`

	From coordinates `json:"from"`
	To   coordinates `json:"to"`

	Price price `json:"price"`

	Status string `json:"status"`
}

var TripExample1 = Trip{
	Id:       "17d52b7c-6f6f-462c-b5e7-98e7fdd14ca9",
	DriverId: "c5ced280-dd79-411f-b699-bb1ef010cd77",
	From:     coordinates{10.0, 10.0},
	To:       coordinates{12.0, 12.0},
	Price:    price{decimal.New(99, 5), "RUB"},
	Status:   "DRIVER_SEARCH",
}

var TripExample2 = Trip{
	Id:       "34ecca26-f17d-491f-9d5d-ef4db5c60876",
	DriverId: "23613a20-5787-42ef-ab81-3524a8e0c33f",
	From:     coordinates{44.0, 38.0},
	To:       coordinates{45.0, 37.0},
	Price:    price{decimal.New(1235, 51), "RUB"},
	Status:   "DRIVER_SEARCH",
}
