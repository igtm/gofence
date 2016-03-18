// Package geofence provides multiple algorithms for use in geofencing
// leverages the diglet go library
package geofence

import (
	"encoding/json"
	"fmt"
	"github.com/buckhx/diglet/geo"
)

const (
	RtreeFence       = "rtree"
	QuadRtreeFence   = "qrtree"
	QuadTreeFence    = "qtree"
	BruteForceFence  = "brute"
	BoundingBoxFence = "bbox"
	CityBruteFence   = "city"
	CityBoxFence     = "city-bbox"
)

// Just a list of the fence types
var FenceLabels = []string{
	RtreeFence, BruteForceFence, QuadTreeFence,
	QuadRtreeFence, CityBruteFence, CityBoxFence,
	BoundingBoxFence,
}

// Interface for algortithms to implement.
type GeoFence interface {
	// Indexes this feature
	Add(f *geo.Feature)
	// Get all features that contain this coordinate
	Get(c geo.Coordinate) []*geo.Feature
}

// Get the rtree geofence as a default. This is the most flexible and will meet most cases
func NewFence() GeoFence {
	return NewRfence()
}

// label is a string from FenceLabels
// Zoom only applies to q-based fences
func GetFence(label string, zoom int) (fence GeoFence, err error) {
	switch label {
	case RtreeFence:
		fence = NewRfence()
	case BruteForceFence:
		fence = NewBruteFence()
	case QuadTreeFence:
		fence = NewQfence(zoom)
	case QuadRtreeFence:
		fence = NewQrfence(zoom)
	case BoundingBoxFence:
		fence = NewBboxFence()
	case CityBruteFence:
		fence = NewCityFence()
	case CityBoxFence:
		fence = NewCityBboxFence()
	default:
		err = fmt.Errorf("Bad fence type: %s", label)
	}
	return
}

// Searchs the fence for the query string. Query should be a geojson point feature.
// Returns a string of the query with a property key 'fences' which is a list of
// the property object of the features that contain the query.
func GeojsonSearch(fence GeoFence, query []byte) (result []byte, err error) {
	gq, err := geo.UnmarshalGeojsonFeature(string(query))
	if err != nil {
		return
	}
	q, err := geo.GeojsonFeatureAdapter(gq)
	if err != nil {
		return
	}
	c := q.Geometry[0].Head() // it's a point
	matchs := fence.Get(c)    // it's a point
	fences := make([]map[string]interface{}, len(matchs))
	for i, match := range matchs {
		fences[i] = match.Properties
	}
	if gq.Properties == nil {
		gq.Properties = make(map[string]interface{}, 1)
	}
	gq.Properties["fences"] = fences
	result, err = json.Marshal(gq)
	return
}
