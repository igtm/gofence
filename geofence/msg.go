package geofence

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/igtm/diglet/geo"
)

//go:generate ffjson msg.go

// TODO change interface{} -> json.RawMessage
type Properties map[string]interface{}

type PointMessage struct {
	Type       string        `json:"type"`
	Properties Properties    `json:"properties"`
	Geometry   PointGeometry `json:"geometry"`
}

type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type ResponseMessage struct {
	Query  PointMessage `json:"query"`
	Fences []Properties `json:"fences"`
}

func newPoint(c geo.Coordinate, props map[string]interface{}) *PointMessage {
	return &PointMessage{
		Type:       "Feature",
		Properties: Properties(props),
		Geometry: PointGeometry{
			Type:        "Point",
			Coordinates: []float64{c.Lon, c.Lat}, // flip
		},
	}
}

func UnmarshalPoint(raw []byte) (point *PointMessage, err error) {
	err = json.Unmarshal(raw, &point)
	if err == nil {
		if point == nil || point.Type != "Feature" || point.Geometry.Type != "Point" {
			err = errors.New("Invalid UnmarshalPoint")
		}
	}
	return
}

// Writes a msg using json encoding
func WriteJson(w io.Writer, msg interface{}) (err error) {
	buf, err := json.Marshal(&msg)
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	//ffjson.Pool(buf)
	return
}

// Searchs the fence for the contents of body. Body should be a geojson point feature.
// Returns a json string of the body with a property key 'fences' which is a list of
// the property object of the features that contain the query.
func GeojsonSearch(idx FenceIndex, name string, body []byte) (result ResponseMessage, err error) {
	point, err := UnmarshalPoint(body)
	if err != nil {
		return
	}
	c := geo.Coordinate{point.Geometry.Coordinates[1], point.Geometry.Coordinates[0]} //geojson swap
	matchs, err := idx.Search(name, c)
	if err != nil {
		return
	}
	fences := make([]Properties, len(matchs))
	for i, fence := range matchs {
		fences[i] = fence.Properties
	}
	result = ResponseMessage{
		Query:  *point,
		Fences: fences,
	}
	return
}
