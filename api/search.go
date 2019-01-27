package api

import (
	"encoding/json"
	"github/ananto/driver_search/storages"
	"log"
	"net/http"
)

// SearchQuery is the struct of a search query
// Distance is in meters and it's a search radius
// Limit is the maximum amount of results to be displayed
type SearchQuery struct {
	Limit    int     `json:"limit"`
	Distance float64 `json:"distance"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

// DriverSearch searches driver from a specific lat lon
func DriverSearch(w http.ResponseWriter, r *http.Request) {
	/*
		WITHDIST: Also return the distance of the returned items from    the specified center. The distance is returned in the same unit as the unit specified as the radius argument of the command.

		WITHCOORD: Also return the longitude,latitude coordinates of the  matching items.

		WITHHASH: Also return the raw geohash-encoded sorted set score of the item, in the form of a 52 bit unsigned integer. This is only useful for low level hacks or debugging and is otherwise of little interest for the general user.
	*/
	rClient := storages.GetRedisClient()
	var search SearchQuery
	if err := json.NewDecoder(r.Body).Decode(&search); err != nil {
		log.Printf("could not decode request %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	drivers := rClient.SearchDrivers(search.Limit, search.Lat, search.Lon, search.Distance)
	data, err := json.Marshal(drivers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}
