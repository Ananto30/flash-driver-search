package events

import (
	"github.com/Ananto30/govent"
)



var publisher chan govent.EventObject

// NewPublisher creates a new publisher
func NewPublisher() {
	publisher = govent.NewEventPublisher()
}

// GetPublisher returns the default publisher
func GetPublisher() chan govent.EventObject {
	return publisher
}

// my event types, registering to govent
const (
	Search govent.EventType = iota
	Found  govent.EventType = iota
)

// SearchEvent is an event type for searching
type SearchEvent struct {
	ID       string
	UserID   string
	Lat, Lon float64
	DriverID string
}

// FoundEvent is an event type for found
type FoundEvent struct {
	ID       string
	UserID   string
	Lat, Lon float64
	DriverID string
}