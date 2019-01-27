package tasks

import (
	"errors"
	"fmt"
	"github/ananto/driver_search/storages"
	"log"
	"strconv"
	"time"
)

// Request invalid reasons
var (
	ErrExpired  = errors.New("request expired")
	ErrCanceled = errors.New("request canceled")
)

// RequestDriverTask is a simple struct that contains information about user, request and driver
type RequestDriverTask struct {
	ID       string
	UserID   string
	Lat, Lon float64
	DriverID string
}

// NewRequestDriverTask creates a new request
func NewRequestDriverTask(id, userID string, lat, lon float64) *RequestDriverTask {
	return &RequestDriverTask{
		ID:     id,
		UserID: userID,
		Lat:    lat,
		Lon:    lon,
	}
}

// Run is the function for executing the task, this task validating the request and launches another goroutine called 'doSearch' which does the search.
func (r *RequestDriverTask) Run(notifier chan bool) {
	ticker := time.NewTicker(time.Second * 5)

	// With the done channel, we receive if the driver was found
	done := make(chan struct{})
	for {
		// The select statement lets a goroutine wait on multiple communication operations.
		select {
		case <-ticker.C:
			err := r.validateRequest()
			switch err {
			case nil:
				log.Println(fmt.Sprintf("Search Driver - Request %s for Lat: %f and Lon: %f", r.ID, r.Lat, r.Lon))
				go r.doSearch(done)
			case ErrExpired:
				sendInfo(r, "Sorry, we didn't find any driver.")
				notifier <- false
				return
			case ErrCanceled:
				log.Printf("Request %s has been canceled.", r.ID)
				return
			default:
				log.Printf("Unexpected error; %v", err)
				return
			}
		case _, notDone := <-done:
			if !notDone {
				sendInfo(r, fmt.Sprintf("Driver %s found", r.DriverID))
				notifier <- true
				ticker.Stop()
				return
			}
		}
	}
}

// validateRequest validates if the request is valid and return an error like a reason in case not.
func (r *RequestDriverTask) validateRequest() error {
	rClient := storages.GetRedisClient()
	keyValue, err := rClient.Get(r.ID).Result()
	if err != nil {
		return ErrExpired
	}
	isActive, _ := strconv.ParseBool(keyValue)
	if !isActive {
		return ErrCanceled
	}
	return nil
}

// doSearch do search of driver and close to the channel.
func (r *RequestDriverTask) doSearch(done chan struct{}) {
	rClient := storages.GetRedisClient()
	drivers := rClient.SearchDrivers(1, r.Lat, r.Lon, 2000)
	if len(drivers) == 1 {
		// Driver found
		// Remove driver location from redis
		rClient.RemoveDriverLocation(drivers[0].Name)
		r.DriverID = drivers[0].Name
		close(done)
	}
	return
}

// sendInfo is a demo method to send information to user, we can use push or socket to notify user
func sendInfo(r *RequestDriverTask, message string) {
	log.Println("Message to user: ", r.UserID)
	log.Println(message)
}
