package geofence

import (
	"github.com/buckhx/diglet/geo"
	"os"
)

type CityFence struct {
	features map[string][]*geo.Feature
	boros    []*geo.Feature
}

// This requires the NYC_BOROS_PATH envvar to be set to the Borrough Boundaries geojson file
// It can be found here http://www1.nyc.gov/site/planning/data-maps/open-data/districts-download-metadata.page
func NewCityFence() *CityFence {
	path := os.Getenv("NYC_BOROS_PATH")
	if path == "" {
		panic("Missing NYC_BOROS_PATH envvar")
	}
	bfeatures, err := geo.NewGeojsonSource(path, nil).Publish()
	if err != nil {
		panic(err)
	}
	var boros []*geo.Feature
	for b := range bfeatures {
		boros = append(boros, b)
	}
	return &CityFence{
		boros:    boros,
		features: make(map[string][]*geo.Feature, 5),
	}
}

// Features must contain a tag BoroName to match to a burrough
func (u *CityFence) Add(f *geo.Feature) {
	u.features[f.Tags("BoroName")] = append(u.features[f.Tags("BoroName")], f)
}

func (u *CityFence) Get(c geo.Coordinate) []*geo.Feature {
	var bn string
	for _, boro := range u.boros {
		if boro.Contains(c) {
			bn = boro.Tags("BoroName")
			break
		}
	}
	if bn == "" {
		return nil
	}
	var ins []*geo.Feature
	for _, f := range u.features[bn] {
		for _, shp := range f.Geometry {
			if shp.Contains(c) {
				ins = append(ins, f)
			}
		}
	}
	return ins
}
