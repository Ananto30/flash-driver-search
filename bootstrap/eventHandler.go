package bootstrap

import (
	"github/ananto/driver_search/events"
	"github/ananto/driver_search/tasks"

	"github.com/Ananto30/govent"
)

func LoadEventHandler() {
	govent.Subscribe(tasks.RunSearchEvent, events.Search)
}

func LoadNewPublisher() {
	events.NewPublisher()
}