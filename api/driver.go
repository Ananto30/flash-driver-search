package api

import (
	"encoding/json"
	"github/ananto/driver_search/storages"
	"log"
	"net/http"
)

// Driver is the structure of driver
type Driver struct {
	DriverID string  `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

// DriverTrack saves driver's location in redis
// Other databases can also be introduced
func DriverTrack(w http.ResponseWriter, r *http.Request) {
	rClient := storages.GetRedisClient()
	var driver Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		log.Printf("could not decode request %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	rClient.AddDriverLocation(driver.Lon, driver.Lat, driver.DriverID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"success\"}"))
	return
}

// AllDrivers shows all the drivers in redis
// func AllDrivers(w http.ResponseWriter, r *http.Request) {
// 	rClient := storages.GetRedisClient()
// 	drivers := rClient.AllDrivers()
// 	data, err := json.Marshal(drivers)
// 	fmt.Print(data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// 	return

// }
