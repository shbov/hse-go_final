package event_type

type EventType string

const (
	CREATED  EventType = "trip.event.created"
	ACCEPTED EventType = "trip.event.accepted"
	CANCELED EventType = "trip.event.cancelled"
	ENDED    EventType = "trip.event.ended"
	STARTED  EventType = "trip.event.started"
)
