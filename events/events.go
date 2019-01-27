package events

import "github.com/Ananto30/govent"

// my event types, registering to govent
const (
	Search govent.EventType = iota
	Found  govent.EventType = iota
)

// MessageEvent is an event type for messaging
type SearchEvent struct {
	ID       string
	UserID   string
	Lat, Lon float64
	DriverID string
}

var publisher chan govent.EventObject

func GetPublisher() chan govent.EventObject {
	publisher := govent.NewEventPublisher()
	return publisher
}
