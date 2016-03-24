package geofence_test

import (
	"bytes"
	"testing"

	"github.com/buckhx/gofence/geofence"
)

func TestMessageRoundTrip(t *testing.T) {
	tests := []struct {
		msg     string
		success bool
	}{
		{"{\"type\":\"Feature\",\"properties\":{},\"geometry\":{\"type\":\"Point\",\"coordinates\":[10,10]}}", true},
		{"{\"type\":\"Feature\",\"properties\":null,\"geometry\":{\"type\":\"Point\",\"coordinates\":[10,10]}}", true},
		{"null", false},
		{"{}", false},
		{"{\"type\":\"Feature\",\"properties\":{},\"geometry\":{\"type\":\"Polygon\",\"coordinates\":[10,10]}}", false},
		{"{\"type\":\"Feature\",\"properties\":{,\"geometry\":{\"type\":\"Point\",\"coordinates\":[10,10]}}", false},
	}
	for i, test := range tests {
		p, err := geofence.UnmarshalPoint([]byte(test.msg))
		if (err == nil) != test.success {
			t.Errorf("Invalid UnmarshalPoint() %d %q - %s", i, test.msg, err)
		} else if err != nil {
			return
		}
		w := bytes.NewBuffer(nil)
		err = geofence.WriteJson(w, p)
		if err != nil {
			t.Errorf("Error in WriteJson %d %q - %s", i, test.msg, err)
		}
		msg := w.String()
		if msg != test.msg {
			t.Errorf("Invalid WriteJson %d %q -> %q", i, test.msg, msg)
		}
	}
}
