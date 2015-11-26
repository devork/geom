package ewkb

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePoint(t *testing.T) {
	// {"00e000000100006c343ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 27700, XYZMS, POINT, []float64{1, 1, 1, 1}}, // EWKB XDR
	// {"01010000e0346c0000000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 27700, XYZMS, POINT, []float64{1, 1, 1, 1}}, // EWKB HDR
	// {"0000000bb93ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 0, XYZM, POINT, []float64{1, 1, 1, 1}},              // WKB XDR
	// {"01b90b0000000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 0, XYZM, POINT, []float64{1, 1, 1, 1}},

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
