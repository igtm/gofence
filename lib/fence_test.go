package geofence_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/buckhx/diglet/geo"
	"github.com/buckhx/gofence/lib"
)

var ues *geo.Feature
var tracts, result []*geo.Feature
var museums = map[string]geo.Coordinate{
	"guggenheim":      {40.7830, -73.9590},
	"met":             {40.7788, -73.9621},
	"moma":            {40.7615, -73.9777},
	"whitney":         {40.7396, -74.0089},
	"old whitney":     {40.7732, -73.9641},
	"natural history": {40.7806, -73.9747},
	"brooklyn":        {40.6713, -73.9638},
	"louvre":          {48.8611, 2.3364},
}

func TestFences(t *testing.T) {
	tests := []struct {
		museum   string
		contains bool
	}{
		{"guggenheim", true},
		{"met", true},
		{"old whitney", true},
		{"whitney", false},
		{"moma", false},
		{"natural history", false},
		{"brooklyn", false},
		{"louvre", false},
	}
	idx := geofence.NewFenceIndex()
	for _, fn := range geofence.FenceLabels {
		fence, err := geofence.GetFence(fn, 10)
		if err != nil {
			// City fences need NYC_BOROS_PATH and we don't always want to test them
			t.Logf("Skipping %q because - %s", fn, err)
			continue
		}
		idx.Set(fn, fence)
		fence.Add(ues)
		for _, test := range tests {
			// Search test
			c := museums[test.museum]
			if (len(fence.Get(c)) == 0) == test.contains {
				t.Errorf("Invalid search %q %q %s", fn, test.museum, c)
			}
			// Index test
			if matchs, err := idx.Search(fn, c); err != nil {
				t.Errorf("Error index search %q - $s", fn, err)
			} else if (len(matchs) == 0) == test.contains {
				t.Errorf("Invalid index search %q %q %s", fn, test.museum, c)
			}
			// Encoding test
			p := &geofence.PointMessage{
				Type:       "Feature",
				Properties: geofence.Properties{"name": []byte(test.museum)}, //TODO fix this
				Geometry:   geofence.PointGeometry{Type: "Point", Coordinates: []float64{c.Lon, c.Lat}},
			}
			b := bytes.NewBuffer(nil)
			err = geofence.WriteJson(b, p)
			if err != nil {
				t.Errorf("Error writing json %s", err)
			}
			res, err := geofence.GeojsonSearch(idx, fn, b.Bytes())
			if err != nil {
				t.Errorf("Error GeojsonSearch %s", err)
			}
			if (len(res.Fences) == 0) == test.contains {
				t.Errorf("Invalid GeojsonSearch %q %q %s", fn, test.museum, c)
			}
		}
	}
}

func BenchmarkBrute(b *testing.B) {
	fence := geofence.NewBruteFence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func BenchmarkCity(b *testing.B) {
	fence, err := geofence.NewCityFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityFence' because %s", err)
		return
	}
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func BenchmarkBbox(b *testing.B) {
	fence := geofence.NewBboxFence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func BenchmarkCityBbox(b *testing.B) {
	fence, err := geofence.NewCityBboxFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityBboxFence' because %s", err)
		return
	}
	if err != nil {
		return
	}
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func BenchmarkQfence(b *testing.B) {
	fence := geofence.NewQfence(14)
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func BenchmarkRfence(b *testing.B) {
	fence := geofence.NewRfence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["met"])
	}
}

func TestMain(m *testing.M) {
	for _, arg := range os.Args {
		// only load tracts if benching
		if strings.Contains(arg, "bench") {
			path := os.Getenv("NYC_TRACTS_PATH")
			if path == "" {
				panic("Missing NYC_TRACTS_PATH envvar")
			}
			features, err := loadGeojson(path)
			if err != nil {
				panic(err)
			}
			tracts = features
			break
		}
	}
	ues = getUpperEastSide()
	os.Exit(m.Run())
}

func loadGeojson(path string) (features []*geo.Feature, err error) {
	source, err := geo.NewGeojsonSource(path, nil).Publish()
	if err != nil {
		return
	}
	for feature := range source {
		features = append(features, feature)
	}
	return
}

func getUpperEastSide() (ues *geo.Feature) {
	shp := geo.NewShape()
	for _, p := range [][]float64{
		{-73.9493, 40.7852}, // w
		{-73.9665, 40.7615}, // s
		{-73.9730, 40.7642}, // e
		{-73.9557, 40.7879}, // n
		{-73.9493, 40.7852}, // w
	} {
		c := geo.Coordinate{p[1], p[0]} //swapped
		shp.Add(c)
	}
	ues = geo.NewPolygonFeature(shp)
	ues.Properties = map[string]interface{}{"BoroName": "Manhattan", "NTAName": "Upper East Side"} // for city
	return
}
