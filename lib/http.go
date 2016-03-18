package geofence

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var fence GeoFence

func ListenAndServe(addr string, gf GeoFence) error {
	fence = gf
	http.HandleFunc("/fence/search", httpSearch)
	log.Printf("Fencing on address %s\n", addr)
	defer log.Printf("Done Fencing\n")
	return http.ListenAndServe(addr, nil)
}

func httpSearch(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB max
	if err != nil {
		http.Error(w, "Body 1 MB max", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "Error closing body", http.StatusInternalServerError)
		return
	}
	result, err := GeojsonSearch(fence, body)
	if err != nil {
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
