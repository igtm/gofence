package geofence

import (
	"errors"
	"os"

	"github.com/buckhx/diglet/geo"
)

// Only for demonstrative purposes
// Checks the containing city first for inclusion, then features. Checks are BoundingBox, then geometry
// This requires the NYC_BOROS_PATH envvar to be set to the Borrough Boundaries geojson file
// It can be found here http://www1.nyc.gov/site/planning/data-maps/open-data/districts-download-metadata.page
type CityBboxFence struct {
	features map[string][]*box
	boros    []*box
}

func NewCityBboxFence() (fence *CityBboxFence, err error) {
	path := os.Getenv("NYC_BOROS_PATH")
	if path == "" {
		err = errors.New("Missing NYC_BOROS_PATH envvar")
		return
	}
	bfeatures, err := geo.NewGeojsonSource(path, nil).Publish()
	if err != nil {
		return
	}
	var boros []*box
	for b := range bfeatures {
		for _, shp := range b.Geometry {
			box := &box{b: shp.BoundingBox(), f: b}
			boros = append(boros, box)
		}
	}
	fence = &CityBboxFence{
		boros:    boros,
		features: make(map[string][]*box, 5),
	}
	return
}

// Features must contain a tag BoroName to match to a burrough
func (u *CityBboxFence) Add(f *geo.Feature) {
	boro, _ := u.features[f.Tags("BoroName")]
	for _, shp := range f.Geometry {
		box := &box{b: shp.BoundingBox(), f: f}
		boro = append(boro, box)
	}
	u.features[f.Tags("BoroName")] = boro
}

func (u *CityBboxFence) Get(c geo.Coordinate) []*geo.Feature {
	var bn string
	for _, boro := range u.boros {
		if boro.b.Contains(c) && boro.f.Contains(c) {
			bn = boro.f.Tags("BoroName")
			break
		}
	}
	if bn == "" {
		return nil
	}
	var ins []*geo.Feature
	for _, box := range u.features[bn] {
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
