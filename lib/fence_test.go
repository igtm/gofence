package geofence_test

import (
	"github.com/buckhx/diglet/geo"
	"github.com/buckhx/gofence/lib"
	"testing"
)

func TestFences(t *testing.T) {
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
	ues := geo.NewPolygonFeature(shp)
	ues.Properties = map[string]interface{}{"BoroName": "Manhattan", "NTAName": "Upper East Side"} // for city
	tests := []struct {
		c        geo.Coordinate
		contains bool
	}{
		{geo.Coordinate{40.7615, -73.9665}, false}, // s Contains != point
		{geo.Coordinate{40.7830, -73.9590}, true},  // guggenheim
		{geo.Coordinate{40.7878, -73.9557}, true},  // n inside
		{geo.Coordinate{40.7484, -73.9857}, false}, // esb
		{geo.Coordinate{-40.7830, 73.9590}, false}, // negative guggenheim
		{geo.Coordinate{40.7889, -73.9557}, false}, // n outside
	}
	for _, fn := range geofence.FenceLabels {
		fence, err := geofence.GetFence(fn, 10)
		if err != nil {
			t.Errorf("Bad GetFence(%s) - %s", fn, err.Error())
		}
		fence.Add(ues)
		for _, test := range tests {
			if (len(fence.Get(test.c)) == 0) == test.contains {
				t.Errorf("Bad containment %q %q", fn, test.c)
			}
		}
	}
}
