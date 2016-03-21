package geofence

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/buckhx/diglet/geo"
	"github.com/julienschmidt/httprouter"
)

var fences FenceIndex

func ListenAndServe(addr string, idx FenceIndex) error {
	log.Printf("Fencing on address %s\n", addr)
	defer log.Printf("Done Fencing\n")
	fences = idx
	//http.Handle("/", fs)
	router := httprouter.New()
	router.GET("/engarde", handleEngarde)
	router.GET("/fence", handleList)
	router.POST("/fence/:name/add", handleAdd)
	router.POST("/fence/:name/search", handleSearch)
	return http.ListenAndServe(addr, router)
}

func respond(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Server", "gofence")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	WriteJson(w, res)
}

func handleList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	WriteJson(w, fences.Keys())
}

func handleAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<26)) // 64 MB max
	if err != nil {
		http.Error(w, "Body 64 MB max", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "Error closing body", http.StatusInternalServerError)
		return
	}
	name := params.ByName("name")
	g, err := geo.UnmarshalGeojsonFeature(string(body))
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	feature, err := geo.GeojsonFeatureAdapter(g)
	if err != nil {
		http.Error(w, "Unable to read geojson feature", http.StatusBadRequest)
		return
	}
	if err := fences.Add(name, feature); err != nil {
		http.Error(w, "Error adding feature "+err.Error(), http.StatusBadRequest)
	}
	respond(w, "success")
}

func handleSearch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB max
	if err != nil {
		http.Error(w, "Body 1 MB max", http.StatusRequestEntityTooLarge)
		return
	}
	if err := r.Body.Close(); err != nil {
		http.Error(w, "Error closing body", http.StatusInternalServerError)
		return
	}
	name := params.ByName("name")
	result, err := GeojsonSearch(fences, name, body)
	if err != nil {
		http.Error(w, "Invalid query "+err.Error(), http.StatusBadRequest)
		return
	}
	respond(w, result)
}

func handleEngarde(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	response := "Touché!"
	w.Header().Set("Server", "gofence")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(response)))
	fmt.Fprint(w, response)
}
