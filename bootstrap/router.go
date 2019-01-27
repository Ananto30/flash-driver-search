package bootstrap

import (
	"github/ananto/driver_search/api"
	"github/ananto/driver_search/api/v2"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func LoadRouter() {
	router := mux.NewRouter()

	// router.HandleFunc("/api/v1/drivers", api.AllDrivers).Methods("GET")
	router.HandleFunc("/api/v1/driver/track", api.DriverTrack).Methods("POST")
	router.HandleFunc("/api/v1/driver/search", api.DriverSearch).Methods("POST")

	// V2
	router.HandleFunc("/api/v2/driver/cancel", v2.CancelRequest).Methods("POST")
	router.HandleFunc("/api/v2/driver/search", v2.SearchV2).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
