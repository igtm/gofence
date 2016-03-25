package geofence_test

import (
	"testing"

	"github.com/buckhx/gofence/geofence"
	"github.com/golang/geo/s2"
)

//TODO use the test structures from fence test

// These tests are really more just making sure the Region interface has been implemented
func TestRegionLoop(t *testing.T) {

	loop := getUpperEastSideLoop()
	degrees := [][]float64{
		{40.7830, -73.9590},
		{40.7788, -73.9621},
		{40.7615, -73.9777},
		{40.7396, -74.0089},
		{40.7732, -73.9641},
		{40.7806, -73.9747},
		{40.6713, -73.9638},
		{48.8611, 2.3364},
	}
	cids := make([]s2.CellID, len(degrees))
	for i, d := range degrees {
		ll := s2.LatLngFromDegrees(d[0], d[1])
		cid := s2.CellIDFromLatLng(ll)
		cids[i] = cid
	}
	/*
		if s2.LatLngFromPoint(loop.CapBound().Center()) != s2.LatLngFromDegrees(40.7852000, -73.9493000) {
			t.Error("Loop cap != 60.0")
		}
		if loop.RectBound().Area() != 0.8434068951180791 {
			t.Error("Loop cap != 0.84")
		}
	*/
	gugg := cids[0]
	louv := cids[7]
	if !loop.ContainsCell(s2.CellFromCellID(gugg)) {
		t.Error("Loop didn't contain gugg cell")
	}
	if !loop.IntersectsCell(s2.CellFromCellID(gugg.Parent(15))) {
		t.Error("Loop didn't intersect gugg cell parent15")
	}
	if !loop.IntersectsCell(s2.CellFromCellID(gugg)) {
		t.Error("Loop didn't intersect gugg cell")
	}
	if loop.ContainsCell(s2.CellFromCellID(louv)) {
		t.Error("Loop contained louv cell")
	}
	if loop.IntersectsCell(s2.CellFromCellID(gugg.Parent(10))) {
		t.Error("Loop intersects louv cell parent10")
	}
}

func getUpperEastSideLoop() *geofence.LoopRegion {
	points := make([]s2.Point, 5)
	for i, d := range [][]float64{
		{-73.9493, 40.7852}, // w
		{-73.9665, 40.7615}, // s
		{-73.9730, 40.7642}, // e
		{-73.9557, 40.7879}, // n
	} {
		ll := s2.LatLngFromDegrees(d[1], d[0]) //.Normalized() //swapped
		p := s2.PointFromLatLng(ll)
		points[i] = p
	}
	return geofence.LoopRegionFromPoints(points)
}
