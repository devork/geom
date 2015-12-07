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

package geom

type Dimension uint32

const (
	XY Dimension = iota
	XYM
	XYZ
	XYZM
	UNKNOWN
)

func (d Dimension) String() string {
	switch d {
	case XY:
		return "XY"
	case XYM:
		return "XYM"
	case XYZ:
		return "XYZ"
	case XYZM:
		return "XYZM"
	default:
		return "UNKNOWN"
	}
}

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

// Geometry interface
type Geometry interface {
	Dimension() Dimension
	SRID() uint32
}

// ----------------------------------------------------------------------------
// Hdr
// ----------------------------------------------------------------------------

// Hdr represents core information about the geometry
type Hdr struct {
	Dim  Dimension
	Srid uint32
}

func (h *Hdr) Dimension() Dimension {
	return h.Dim
}

func (h *Hdr) SRID() uint32 {
	return h.Srid
}

// ----------------------------------------------------------------------------
// Point
// ----------------------------------------------------------------------------

// Point
type Point struct {
	Hdr
	Coordinate
}

// ----------------------------------------------------------------------------
// MultiPoint
// ----------------------------------------------------------------------------

type MultiPoint struct {
	Hdr
	Points []Point
}

// ----------------------------------------------------------------------------
// LineString
// ----------------------------------------------------------------------------

type LineString struct {
	Hdr
	Coordinates []Coordinate
}

// ----------------------------------------------------------------------------
// MultiLineString
// ----------------------------------------------------------------------------

type MultiLineString struct {
	Hdr
	LineStrings []LineString
}

// func (p *MultiLineString) GeoJSON(crs, bbox bool) string {
// 	var sb bytes.Buffer
// 	sb.WriteString(`{"type":"MultiLineString", "coordinates":[`)
//
// 	llimit := len(p.LineStrings) - 1
// 	for lidx, linestring := range p.LineStrings {
// 		sb.WriteString("[")
//
// 		limit := len(linestring.Coordinates) - 1
// 		for idx, coord := range linestring.Coordinates {
// 			coord.appendGeoJSON(&sb)
//
// 			if idx < limit {
// 				sb.WriteString(",")
// 			}
// 		}
//
// 		sb.WriteString("]")
//
// 		if lidx < llimit {
// 			sb.WriteString(",")
// 		}
// 	}
//
// 	sb.WriteString(`],`)
// 	sb.WriteString(`"dim":"`)
// 	sb.WriteString(p.Hdr.dim.String())
// 	sb.WriteString(`"`)
// 	if crs {
// 		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
// 		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 		sb.WriteString(`"}}`)
// 	}
// 	sb.WriteString(`}`)
//
// 	return sb.String()
// }
//
// func (p *MultiLineString) EWKT() string {
// 	var sb bytes.Buffer
//
// 	sb.WriteString("SRID=")
// 	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 	sb.WriteString(";MULTILINESTRING()")
//
// 	return sb.String()
// }
//
// func (p *MultiLineString) String() string {
// 	return p.EWKT()
// }

// ----------------------------------------------------------------------------
// Polygon
// ----------------------------------------------------------------------------

type Polygon struct {
	Hdr
	Rings []LinearRing
}

// ----------------------------------------------------------------------------
// MultiPolygon
// ----------------------------------------------------------------------------

type MultiPolygon struct {
	Hdr
	Polygons []Polygon
}

// func (p *MultiPolygon) GeoJSON(crs, bbox bool) string {
// 	var sb bytes.Buffer
// 	sb.WriteString(`{"type":"MultiPolygon", "coordinates":[`)
//
// 	plimit := len(p.Polygons) - 1
// 	for pidx, polygon := range p.Polygons {
// 		sb.WriteString("[")
//
// 		rlimit := len(polygon.Rings) - 1
// 		for ridx, lring := range polygon.Rings {
// 			lring.appendGeoJSON(&sb)
//
// 			if ridx < rlimit {
// 				sb.WriteString(",")
// 			}
// 		}
//
// 		sb.WriteString("]")
//
// 		if pidx < plimit {
// 			sb.WriteString(",")
// 		}
// 	}
//
// 	sb.WriteString(`],`)
// 	sb.WriteString(`"dim":"`)
// 	sb.WriteString(p.Hdr.dim.String())
// 	sb.WriteString(`"`)
// 	if crs {
// 		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
// 		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 		sb.WriteString(`"}}`)
// 	}
// 	sb.WriteString(`}`)
// 	return sb.String()
// }
//
// func (p *MultiPolygon) EWKT() string {
// 	var sb bytes.Buffer
//
// 	sb.WriteString("SRID=")
// 	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 	sb.WriteString(";MULTIPOLYGON()")
//
// 	return sb.String()
// }
//
// func (p *MultiPolygon) String() string {
// 	return p.EWKT()
// }

// ----------------------------------------------------------------------------
// GeometryCollection (a misnomer IMHO - should be called MultiGeometry)
// ----------------------------------------------------------------------------

type GeometryCollection struct {
	Hdr
	Geometries []Geometry
}

// func (p *GeometryCollection) GeoJSON(crs, bbox bool) string {
// 	var sb bytes.Buffer
// 	sb.WriteString(`{"type":"GeometryCollection", "geometries":[`)
// 	limit := len(p.Geometries) - 1
// 	for idx, geom := range p.Geometries {
// 		sb.WriteString(geom.GeoJSON(false, false))
//
// 		if idx < limit {
// 			sb.WriteString(",")
// 		}
// 	}
// 	sb.WriteString(`],`)
// 	sb.WriteString(`"dim":"`)
// 	sb.WriteString(p.Hdr.dim.String())
// 	sb.WriteString(`"`)
// 	if crs {
// 		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
// 		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 		sb.WriteString(`"}}`)
// 	}
// 	sb.WriteString(`}`)
//
// 	return sb.String()
// }
//
// func (p *GeometryCollection) EWKT() string {
// 	var sb bytes.Buffer
//
// 	sb.WriteString("SRID=")
// 	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
// 	sb.WriteString(";GEOMETRYCOLLECTION()")
//
// 	return sb.String()
// }
//
// func (p *GeometryCollection) String() string {
// 	return p.EWKT()
// }

// ----------------------------------------------------------------------------
// LinearRing
// ----------------------------------------------------------------------------

type LinearRing struct {
	Coordinates []Coordinate
}

// ----------------------------------------------------------------------------
// Coordinate
// ----------------------------------------------------------------------------

// Coordinate type
type Coordinate []float64
