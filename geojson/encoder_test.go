package geojson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/devork/geom"
	"github.com/stretchr/testify/assert"
)

func TestEncodePoint(t *testing.T) {
	expected := &geom.Point{
		geom.Hdr{
			Dim:  geom.XYZ,
			Srid: 27700,
		},
		[]float64{-0.118340, 51.503475, -12.34567890},
	}

	var sb bytes.Buffer
	err := Encode(expected, &sb)

	if err != nil {
		t.Fatalf("failed to marshal point: %s", err)
	}

	got := make(map[string]interface{})
	err = json.Unmarshal(sb.Bytes(), &got)

	if err != nil {
		t.Fatalf("Failed to parse generated GeoJSON: error %s", err)
	}

	t.Logf("%s\n", string(sb.Bytes()))

	assert.Equal(t, "Point", got["type"])

	coords := got["coordinates"].([]interface{})
	assert.InDelta(t, -0.118340, coords[0].(float64), 1e-9)
	assert.InDelta(t, 51.503475, coords[1].(float64), 1e-9)
	assert.InDelta(t, -12.34567890, coords[2].(float64), 1e-9)
}

func TestEncodeLineString(t *testing.T) {
	expected := &geom.LineString{
		geom.Hdr{
			Dim:  geom.XYZ,
			Srid: 27700,
		},
		[]geom.Coordinate{
			[]float64{100.0, 0.0, 1.546},
			[]float64{101.0, 1.0, 2.345},
		},
	}

	var sb bytes.Buffer
	err := Encode(expected, &sb)

	if err != nil {
		t.Fatalf("failed to marshal linestring: %s", err)
	}

	got := make(map[string]interface{})
	err = json.Unmarshal(sb.Bytes(), &got)

	if err != nil {
		t.Fatalf("Failed to parse generated GeoJSON: error %s", err)
	}

	t.Logf("%s\n", string(sb.Bytes()))

	assert.Equal(t, "LineString", got["type"])

	coords := got["coordinates"].([]interface{})

	assert.Equal(t, 2, len(coords))

	coord0 := coords[0].([]interface{})
	assert.InDelta(t, 100, coord0[0].(float64), 1e-9)
	assert.InDelta(t, 0, coord0[1].(float64), 1e-9)
	assert.InDelta(t, 1.546, coord0[2].(float64), 1e-9)

	coord1 := coords[1].([]interface{})
	assert.InDelta(t, 101, coord1[0].(float64), 1e-9)
	assert.InDelta(t, 1, coord1[1].(float64), 1e-9)
	assert.InDelta(t, 2.345, coord1[2].(float64), 1e-9)

}
