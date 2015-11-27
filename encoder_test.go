package ewkb

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePoint(t *testing.T) {
	datasets := []struct {
		data     Geometry
		expected string
	}{
		{&Point{Hdr{XYZMS, 27700, POINT}, Coordinate{1, 1, 1, 1}}, "00e000000100006c343ff00000000000003ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYZM, 0, POINT}, Coordinate{1, 1, 1, 1}}, "0000000bb93ff00000000000003ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYZS, 27700, POINT}, Coordinate{1, 1, 1}}, "00a000000100006c343ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYZ, 0, POINT}, Coordinate{1, 1, 1}}, "00000003e93ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYMS, 27700, POINT}, Coordinate{1, 1, 1}}, "006000000100006c343ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYM, 0, POINT}, Coordinate{1, 1, 1}}, "00000007d13ff00000000000003ff00000000000003ff0000000000000"},
		{&Point{Hdr{XYS, 27700, POINT}, Coordinate{1, 1}}, "002000000100006c343ff00000000000003ff0000000000000"},
		{&Point{Hdr{XY, 0, POINT}, Coordinate{1, 1}}, "00000000013ff00000000000003ff0000000000000"},
	}

	for _, dataset := range datasets {
		var w = new(bytes.Buffer)
		err := Encode(dataset.data, w)

		if err != nil {
			t.Fatalf("Failed to encode Point geometry: err = %s", err)
		}
		data := hex.EncodeToString(w.Bytes())

		assert.Equal(t, dataset.expected, strings.ToLower(data))

	}

}

func TestEncodeLineString(t *testing.T) {
	datasets := []struct {
		data     Geometry
		expected string
	}{
		{&LineString{Hdr{XYS, 27700, LINESTRING}, []Coordinate{{30, 10}, {10, 30}, {40, 40}}}, "002000000200006c3400000003403e00000000000040240000000000004024000000000000403e00000000000040440000000000004044000000000000"},
	}

	for _, dataset := range datasets {
		var w = new(bytes.Buffer)
		err := Encode(dataset.data, w)

		if err != nil {
			t.Fatalf("Failed to encode LineString geometry: err = %s", err)
		}
		data := hex.EncodeToString(w.Bytes())

		assert.Equal(t, dataset.expected, strings.ToLower(data))

	}

}

func TestEncodePolygon(t *testing.T) {
	datasets := []struct {
		data     Geometry
		expected string
	}{
		{
			&Polygon{
				Hdr{XYS, 27700, POLYGON},
				[]LinearRing{{
					[]Coordinate{
						{30, 10},
						{40, 40},
						{20, 40},
						{10, 20},
						{30, 10},
					},
				},
				},
			},
			"002000000300006c340000000100000005403e0000000000004024000000000000404400000000000040440000000000004034000000000000404400000000000040240000000000004034000000000000403e0000000000004024000000000000"},
	}

	for _, dataset := range datasets {
		var w = new(bytes.Buffer)
		err := Encode(dataset.data, w)

		if err != nil {
			t.Fatalf("Failed to encode Polygon geometry: err = %s", err)
		}
		data := hex.EncodeToString(w.Bytes())

		assert.Equal(t, dataset.expected, strings.ToLower(data))

	}

}
