package geofence

import (
	"github.com/buckhx/diglet/geo"
)

type box struct {
	b geo.Box
	f *geo.Feature
}

type BboxFence struct {
	boxes []*box
}

func NewBboxFence() *BboxFence {
	return &BboxFence{}
}

func (b *BboxFence) Add(f *geo.Feature) {
	for _, shp := range f.Geometry {
		box := &box{b: shp.BoundingBox(), f: f}
		b.boxes = append(b.boxes, box)
	}
}

func (b *BboxFence) Get(c geo.Coordinate) []*geo.Feature {
	var ins []*geo.Feature
	for _, box := range b.boxes {
		if box.b.Contains(c) {
			for _, shp := range box.f.Geometry {
				if shp.Contains(c) {
					ins = append(ins, box.f)
				}
			}
		}
	}
	return ins
}
