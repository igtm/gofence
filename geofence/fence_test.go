package geofence_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/buckhx/diglet/geo"
	"github.com/buckhx/gofence/geofence"
)

const (
	TEST_ZOOM = 14
)

var (
	ues            *geo.Feature
	s2f            geofence.GeoFence
	tracts, result []*geo.Feature
	museums        = map[string]geo.Coordinate{
		"guggenheim":      {40.7830, -73.9590},
		"met":             {40.7788, -73.9621},
		"moma":            {40.7615, -73.9777},
		"whitney":         {40.7396, -74.0089},
		"old whitney":     {40.7732, -73.9641},
		"natural history": {40.7806, -73.9747},
		"brooklyn":        {40.6713, -73.9638},
		"louvre":          {48.8611, 2.3364},
	}
)

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
		fence, err := geofence.GetFence(fn, TEST_ZOOM)
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

func BenchmarkBruteGet(b *testing.B) {
	fence := geofence.NewBruteFence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkCityGet(b *testing.B) {
	fence, err := geofence.NewCityFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityFence' because %s", err)
		return
	}
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkBboxGet(b *testing.B) {
	fence := geofence.NewBboxFence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkCityBboxGet(b *testing.B) {
	fence, err := geofence.NewCityBboxFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityBboxFence' because %s", err)
		return
	}
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkQfenceGet(b *testing.B) {
	fence := geofence.NewQfence(TEST_ZOOM)
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkRfenceGet(b *testing.B) {
	fence := geofence.NewRfence()
	for _, tract := range tracts {
		fence.Add(tract)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = fence.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkS2fenceGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// interior @ Z18
		result = s2f.Get(museums["old whitney"])
		if len(result) != 1 {
			b.Fatal("Incorrect Get() result")
		}
	}
}

func BenchmarkBruteAdd(b *testing.B) {
	fence := geofence.NewBruteFence()
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkCityAdd(b *testing.B) {
	fence, err := geofence.NewCityFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityFence' because %s", err)
		return
	}
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkBboxAdd(b *testing.B) {
	fence := geofence.NewBboxFence()
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkCityBboxAdd(b *testing.B) {
	fence, err := geofence.NewCityBboxFence()
	if err != nil {
		fmt.Printf("Skipping benchmark for 'CityBboxFence' because %s", err)
		return
	}
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkQfenceAdd(b *testing.B) {
	fence := geofence.NewQfence(TEST_ZOOM)
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkRfenceAdd(b *testing.B) {
	fence := geofence.NewRfence()
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
	}
}

func BenchmarkS2fenceAdd(b *testing.B) {
	fence := geofence.NewS2fence(TEST_ZOOM)
	for n := 0; n < b.N; n++ {
		tract := tracts[n%len(tracts)]
		fence.Add(tract)
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
			fmt.Println("Loading s2fence...")
			s2f = geofence.NewS2fence(TEST_ZOOM)
			for _, tract := range tracts {
				//fmt.Printf("s2fence adding feature %d\n", i)
				s2f.Add(tract)
			}
			fmt.Println("Loaded s2fence!")
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
