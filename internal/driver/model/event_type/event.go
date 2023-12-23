package event_type

type EventType string

const (
	CREATED  EventType = "event.trip.created"
	ACCEPTED EventType = "event.trip.accepted"
	CANCELED EventType = "event.trip.cancelled"
	ENDED    EventType = "event.trip.ended"
	STARTED  EventType = "event.trip.started"
)
