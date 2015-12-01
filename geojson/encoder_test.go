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
