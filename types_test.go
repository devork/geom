/*
Copyright [2015] Alex Davies-Moore

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ewkb

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_GeoJSON(t *testing.T) {
	expected := &Point{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: POINT,
		},
		[]float64{100.0, 0.0},
	}

	geojson := expected.GeoJSON(true, true)

	got := make(map[string]interface{})
	err := json.Unmarshal([]byte(geojson), &got)

	if err != nil {
		t.Fatalf("Failed to parse generated GeoJSON: error %s", err)
	}

	//t.Logf("%#v\n", got)

	assert.Equal(t, "Point", got["type"])
	assert.Equal(t, "XYZ", got["dim"])

	coords := got["coordinates"].([]interface{})
	assert.InDelta(t, 100, coords[0].(float64), 1e-9)
	assert.InDelta(t, 0, coords[1].(float64), 1e-9)

	crs := got["crs"].(map[string]interface{})
	assert.Equal(t, "name", crs["type"])

	props := crs["properties"].(map[string]interface{})
	assert.Equal(t, "EPSG:27700", props["name"])
}

func TestLineString_GeoJSON(t *testing.T) {
	expected := &LineString{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: LINESTRING,
		},
		[]Coordinate{
			[]float64{100.0, 0.0},
			[]float64{101.0, 1.0},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestPolygon_GeoJSON(t *testing.T) {
	expected := &Polygon{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: POLYGON,
		},
		[]LinearRing{
			LinearRing{
				[]Coordinate{
					[]float64{100.0, 0.0},
					[]float64{101.0, 0.0},
					[]float64{101.0, 1.0},
					[]float64{100.0, 1.0},
					[]float64{100.0, 0.0},
				},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestPolygonHoles_GeoJSON(t *testing.T) {
	expected := &Polygon{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: POLYGON,
		},
		[]LinearRing{
			LinearRing{
				[]Coordinate{
					[]float64{100.0, 0.0},
					[]float64{101.0, 0.0},
					[]float64{101.0, 1.0},
					[]float64{100.0, 1.0},
					[]float64{100.0, 0.0},
				},
			},
			LinearRing{
				[]Coordinate{
					[]float64{100.2, 0.2},
					[]float64{100.8, 0.2},
					[]float64{100.8, 0.8},
					[]float64{100.2, 0.8},
					[]float64{100.2, 0.2},
				},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestMultiPoint_GeoJSON(t *testing.T) {
	expected := &MultiPoint{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: MULTIPOINT,
		},
		[]Point{
			Point{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: POINT,
				},
				[]float64{100.0, 0.0},
			},
			Point{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: POINT,
				},
				[]float64{101.0, 1.0},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestMultiLineString_GeoJSON(t *testing.T) {
	expected := &MultiLineString{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: MULTILINESTRING,
		},
		[]LineString{
			LineString{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: LINESTRING,
				},
				[]Coordinate{
					[]float64{100.0, 0.0},
					[]float64{101.0, 1.0},
				},
			},
			LineString{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: LINESTRING,
				},
				[]Coordinate{
					[]float64{102.0, 2.0},
					[]float64{103.0, 3.0},
				},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestMultiPolygon_GeoJSON(t *testing.T) {
	expected := &MultiPolygon{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: MULTIPOLYGON,
		},
		[]Polygon{
			Polygon{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: POLYGON,
				},
				[]LinearRing{
					LinearRing{
						[]Coordinate{
							[]float64{100.0, 0.0},
							[]float64{101.0, 0.0},
							[]float64{101.0, 1.0},
							[]float64{100.0, 1.0},
							[]float64{100.0, 0.0},
						},
					},
				},
			},
			Polygon{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: POLYGON,
				},
				[]LinearRing{
					LinearRing{
						[]Coordinate{
							[]float64{100.0, 0.0},
							[]float64{101.0, 0.0},
							[]float64{101.0, 1.0},
							[]float64{100.0, 1.0},
							[]float64{100.0, 0.0},
						},
					},
					LinearRing{
						[]Coordinate{
							[]float64{100.2, 0.2},
							[]float64{100.8, 0.2},
							[]float64{100.8, 0.8},
							[]float64{100.2, 0.8},
							[]float64{100.2, 0.2},
						},
					},
				},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}

func TestGeometryCollection_GeoJSON(t *testing.T) {
	expected := &GeometryCollection{
		Hdr{
			dim:   XYZ,
			srid:  27700,
			gtype: LINESTRING,
		},
		[]Geometry{
			&Point{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: POINT,
				},
				[]float64{100.0, 0.0},
			},
			&LineString{
				Hdr{
					dim:   XYZ,
					srid:  27700,
					gtype: LINESTRING,
				},
				[]Coordinate{
					[]float64{100.0, 0.0},
					[]float64{101.0, 1.0},
				},
			},
		},
	}

	geojson := expected.GeoJSON(true, true)

	t.Logf("%s\n", geojson)
}
