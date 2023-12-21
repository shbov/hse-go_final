package commands

import "time"

type AcceptMessage struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            struct {
		TripId   string `json:"trip_id"`
		DriverId string `json:"driver_id"`
	} `json:"data"`
}

type CancelMessage struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            struct {
		TripId string `json:"trip_id"`
		Reason string `json:"reason"`
	} `json:"data"`
}

type StartEndMessage struct {
	Id              string    `json:"id"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	DataContentType string    `json:"datacontenttype"`
	Time            time.Time `json:"time"`
	Data            struct {
		TripId string `json:"trip_id"`
	} `json:"data"`
}
