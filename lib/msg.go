package geofence

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/pquerna/ffjson/ffjson"
)

//go:generate ffjson msg.go

type Properties map[string]json.RawMessage

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

func UnmarshalPoint(raw []byte) (point *PointMessage, err error) {
	err = ffjson.Unmarshal(raw, &point)
	if err == nil {
		if point == nil || point.Type != "Feature" || point.Geometry.Type != "Point" {
			err = errors.New("Invalid UnmarshalPoint")
		}
	}
	return
}

func WriteJson(w io.Writer, msg interface{}) (err error) {
	buf, err := ffjson.Marshal(&msg)
	_, err = w.Write(buf)
	if err != nil {
		return err
	}
	ffjson.Pool(buf)
	return
}
