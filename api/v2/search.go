package v2

import (
	"encoding/json"
	"fmt"
	"github/ananto/driver_search/storages"
	"github/ananto/driver_search/tasks"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SearchResponse struct {
	RequestorID string
	DriverID    string
}

func SearchV2(w http.ResponseWriter, r *http.Request) {
	rClient := storages.GetRedisClient()

	// A unique key for each request in redis
	// With this key we can know if the request is active or canceled
	requestID, err := rClient.Incr("request_id").Result()
	if err != nil {
		return
	}
	key := strconv.Itoa(int(requestID))

	// Set true value for the key and also the expiration time, this expiration time is the duration that has the request to find a driver.
	rClient.Set(key, true, time.Minute*2)

	body := struct {
		Lat, Lon float64
	}{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Could not decode request %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}

	// publisher := events.GetPublisher()
	// publisher <- govent.EventObject{
	// 	EventType: events.Search,
	// 	Event: events.SearchEvent{
	// 		ID:     key,
	// 		UserID: fmt.Sprintf("requestor_%s", key),
	// 		Lat:    body.Lat,
	// 		Lon:    body.Lon,
	// 	},
	// }

	notifier := make(chan bool, 1)
	rTask := tasks.NewRequestDriverTask(key, fmt.Sprintf("requestor_%s", key), body.Lat, body.Lon)
	go rTask.Run(notifier)

	found := <-notifier
	if found {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"driver found": %s}`, key)))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf(`{"driver not found": %s}`, key)))
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf(`{"request_id": %s}`, key)))
}

func CancelRequest(w http.ResponseWriter, r *http.Request) {
	rClient := storages.GetRedisClient()
	body := struct {
		RequestID string `json:"request_id"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Could not decode request %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	rClient.Set(body.RequestID, false, time.Minute*1)
	w.WriteHeader(http.StatusOK)
	return
}
