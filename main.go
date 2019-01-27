package main

import (
	"github/ananto/driver_search/bootstrap"
)

func main() {
	bootstrap.LoadNewPublisher()
	bootstrap.LoadEventHandler()
	bootstrap.LoadRouter()
}
