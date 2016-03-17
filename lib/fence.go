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

var FenceLabels = []string{
	RtreeFence, BruteForceFence, QuadTreeFence,
	QuadRtreeFence, CityBruteFence, CityBoxFence,
	BoundingBoxFence,
}

type GeoFence interface {
	Add(f *geo.Feature)
	Get(c geo.Coordinate) []*geo.Feature
}

func NewFence() GeoFence {
	return NewRfence()
}

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
