package geofence

import (
	"github.com/buckhx/diglet/geo"
	"github.com/golang/geo/s2"
)

type S2fence struct {
	zoom   int
	covers map[s2.CellID][]cover
}

func NewS2fence(zoom int) *S2fence {
	return &S2fence{
		zoom:   zoom,
		covers: make(map[s2.CellID][]cover),
	}
}

func (s *S2fence) Add(f *geo.Feature) {
	// this will give us a flat covering
	// not ideal, but don't want to deal with bitwise prefixes today
	coverer := NewFlatCoverer(s.zoom)
	for _, shp := range f.Geometry {
		if shp.IsClockwise() {
			shp.Reverse() //s2 wants CCW
		}
		points := make([]s2.Point, len(shp.Coordinates))
		for i, c := range shp.Coordinates[:] {
			points[i] = s2.PointFromLatLng(s2.LatLngFromDegrees(c.Lat, c.Lon))
		}
		region := s2.Region(LoopRegionFromPoints(points))
		bounds := coverer.Covering(region)
		if len(bounds) < 1 {
			continue
		}
		interiors := coverer.InteriorCovering(region)
		cov := cover{
			interior: make(map[s2.CellID]bool, len(interiors)),
			feature:  f,
		}
		for _, cid := range interiors {
			cov.interior[cid] = true
		}
		for _, cid := range bounds {
			s.covers[cid] = append(s.covers[cid], cov)
		}
	}
}

func (s *S2fence) Get(c geo.Coordinate) (matchs []*geo.Feature) {
	cid := s2.CellIDFromLatLng(s2.LatLngFromDegrees(c.Lat, c.Lon)).Parent(s.zoom)
	for _, cov := range s.covers[cid] {
		if _, ok := cov.interior[cid]; ok {
			matchs = append(matchs, cov.feature)
		} else if cov.feature.Contains(c) {
			matchs = append(matchs, cov.feature)
		}
	}
	return
}

// has a feature and an interior for face lookips
type cover struct {
	feature *geo.Feature
	// if this map takes up too much space, change value to an empty struct{}
	interior map[s2.CellID]bool
}

// Making s2.Loop implement s2.Region
type LoopRegion struct {
	*s2.Loop
}

func LoopRegionFromPoints(points []s2.Point) *LoopRegion {
	loop := s2.LoopFromPoints(points)
	return &LoopRegion{loop}
}

func (l *LoopRegion) CapBound() s2.Cap {
	return l.RectBound().CapBound()
}

func (l *LoopRegion) ContainsCell(c s2.Cell) bool {
	for i := 0; i < 4; i++ {
		v := c.Vertex(i)
		if !l.ContainsPoint(v) {
			return false
		}
	}
	return true
}

func (l *LoopRegion) IntersectsCell(c s2.Cell) bool {
	for i := 0; i < 4; i++ {
		crosser := s2.NewChainEdgeCrosser(c.Vertex(i), c.Vertex((i+1)%4), l.Vertex(0))
		for _, v := range l.Vertices()[1:] {
			if crosser.EdgeOrVertexChainCrossing(v) {
				return true
			}
		}
		if crosser.EdgeOrVertexChainCrossing(l.Vertex(0)) { //close the loop
			return true
		}
	}
	return l.ContainsCell(c)
}

// Embeds a s2.RegionCover, but does it's own covering
// Pick the deepest level and normalize a cellunion
// The default coverer didn't trim on the boundary...
type FlatCoverer struct {
	*s2.RegionCoverer
}

func NewFlatCoverer(level int) *FlatCoverer {
	return &FlatCoverer{&s2.RegionCoverer{
		MinLevel: level,
		MaxLevel: level,
		LevelMod: 0,
		MaxCells: 1 << 12,
	}}
}

func (c *FlatCoverer) Covering(r s2.Region) s2.CellUnion {
	var cover []s2.CellID
	cids := c.FastCovering(r.CapBound())
	for _, cid := range cids {
		cell := s2.CellFromCellID(cid)
		if r.IntersectsCell(cell) {
			cover = append(cover, cid)
		}
	}
	return s2.CellUnion(cover)
}

func (c *FlatCoverer) CellUnion(r s2.Region) s2.CellUnion {
	cover := c.Covering(r)
	cover.Normalize()
	return cover
}

func (c *FlatCoverer) InteriorCovering(r s2.Region) s2.CellUnion {
	var cover []s2.CellID
	cids := c.FastCovering(r.CapBound())
	for _, cid := range cids {
		cell := s2.CellFromCellID(cid)
		if r.ContainsCell(cell) {
			cover = append(cover, cid)
		}
	}
	return s2.CellUnion(cover)
}

func (c *FlatCoverer) InteriorCellUnion(r s2.Region) s2.CellUnion {
	cover := c.InteriorCovering(r)
	cover.Normalize()
	return cover
}
