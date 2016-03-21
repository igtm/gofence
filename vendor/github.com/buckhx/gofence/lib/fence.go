// Package geofence provides multiple algorithms for use in geofencing
// leverages the diglet go library
package geofence

import (
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
