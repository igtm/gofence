package geofence

import (
	"github.com/igtm/diglet/geo"
)

// Standard rtree with M=50
type Rfence struct {
	rtree *geo.Rtree
}

func NewRfence() *Rfence {
	return &Rfence{
		rtree: geo.NewRtree(),
	}
}

func (r *Rfence) Add(f *geo.Feature) {
	for _, shp := range f.Geometry {
		if len(shp.Coordinates) > 1 { //single coords don't have an area
			r.rtree.Insert(shp, f)
		}
	}
}

func (r *Rfence) Get(c geo.Coordinate) (matchs []*geo.Feature) {
	nodes := r.rtree.Contains(c)
	for _, n := range nodes {
		feature := n.Feature()
		if feature.Contains(c) {
			matchs = append(matchs, feature)
		}
	}
	return
}
