package geofence

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/buckhx/diglet/geo"
	"github.com/julienschmidt/httprouter"
)

var fences FenceIndex

func ListenAndServe(addr string, idx FenceIndex, profile bool) error {
	log.Printf("Fencing on address %s\n", addr)
	defer log.Printf("Done Fencing\n")
	fences = idx
	//http.Handle("/", fs)
	router := httprouter.New()
	router.GET("/engarde", getEngarde)
	router.GET("/fence", getList)
	router.POST("/fence/:name/add", postAdd)
	router.POST("/fence/:name/search", postSearch)
	router.GET("/fence/:name/search", getSearch)
	if profile {
		attachProfiler(router)
	}
	return http.ListenAndServe(addr, router)
}

func respond(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Server", "gofence")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	WriteJson(w, res)
}

func getList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	WriteJson(w, fences.Keys())
}

func postAdd(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func postSearch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

func getSearch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	query := r.URL.Query()
	lat, err := strconv.ParseFloat(query.Get("lat"), 64)
	if err != nil {
		http.Error(w, "Query param 'lat' required as float", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(query.Get("lon"), 64)
	if err != nil {
		http.Error(w, "Query param 'lon' required as float", http.StatusBadRequest)
		return
	}
	query.Del("lat")
	query.Del("lon")
	c := geo.Coordinate{Lat: lat, Lon: lon}
	name := params.ByName("name")
	matchs, err := fences.Search(name, c)
	if err != nil {
		http.Error(w, "Error search fence "+name, http.StatusBadRequest)
		return
	}
	fences := make([]Properties, len(matchs))
	for i, fence := range matchs {
		fences[i] = fence.Properties
	}
	props := make(map[string]interface{}, len(query))
	for k := range query {
		props[k] = query.Get(k)
	}
	result := ResponseMessage{
		Query:  *newPoint(c, props),
		Fences: fences,
	}
	respond(w, result)
}

func getEngarde(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	response := "Touché!"
	w.Header().Set("Server", "gofence")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(response)))
	fmt.Fprint(w, response)
}

func attachProfiler(router *httprouter.Router) {
	router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
	router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
	router.Handler("GET", "/debug/pprof/heap", pprof.Handler("heap"))
	router.Handler("GET", "/debug/pprof/block", pprof.Handler("block"))
	router.Handler("GET", "/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handler("GET", "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
}
