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
	Id       int    `json:"id"`
	DriverId string `json:"driver_id"`

	From coordinates `json:"from"`
	To   coordinates `json:"to"`

	Price price `json:"price"`

	Status string `json:"status"`
}
