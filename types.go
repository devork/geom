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
	"bytes"
	"strconv"
)

// ----------------------------------------------------------------------------
// Dimension
// ----------------------------------------------------------------------------

// Dimension of the geometry
type Dimension uint16

func (d Dimension) String() string {
	switch d {
	case XY:
		return "XY"
	case XYS:
		return "XYS"
	case XYZ:
		return "XYZ"
	case XYZS:
		return "XYZS"
	case XYM:
		return "XYM"
	case XYMS:
		return "XYMS"
	case XYZM:
		return "XYZM"
	case XYZMS:
		return "XYZMS"
	default:
		return "UNKNOWN"
	}
}

// ----------------------------------------------------------------------------
// GeomType
// ----------------------------------------------------------------------------

// GeomType is the bitmask of the geom
type GeomType uint16

func (g GeomType) String() string {
	switch g {
	case GEOMETRY:
		return "GEOMETRY"
	case POINT:
		return "POINT"
	case LINESTRING:
		return "LINESTRING"
	case POLYGON:
		return "POLYGON"
	case MULTIPOINT:
		return "MULTIPOINT"
	case MULTILINESTRING:
		return "MULTILINESTRING"
	case MULTIPOLYGON:
		return "MULTIPOLYGON"
	case GEOMETRYCOLLECTION:
		return "GEOMETRYCOLLECTION"
	// case CIRCULARSTRING:
	// 	return "CIRCULARSTRING"
	// case COMPOUNDCURVE:
	// 	return "COMPOUNDCURVE"
	// case CURVEPOLYGON:
	// 	return "CURVEPOLYGON"
	// case MULTICURVE:
	// 	return "MULTICURVE"
	// case MULTISURFACE:
	// 	return "MULTISURFACE"
	// case CURVE:
	// 	return "CURVE"
	// case SURFACE:
	// 	return "SURFACE"
	// case POLYHEDRALSURFACE:
	// 	return "POLYHEDRALSURFACE "
	// case TIN:
	// 	return "TIN"
	// case TRIANGLE:
	// 	return "TRIANGLE"
	default:
		return "UNKNOWN"
	}
}

// ----------------------------------------------------------------------------
// Geometry
// ----------------------------------------------------------------------------

// Geometry interface
type Geometry interface {
	Srid() uint32
	Dimension() Dimension
	Type() GeomType
	GeoJSON(crs, bbox bool) string
	EWKT() string
}

// ----------------------------------------------------------------------------
// Hdr
// ----------------------------------------------------------------------------

// Hdr represents core information about the geometry
type Hdr struct {
	dim   Dimension
	srid  uint32
	gtype GeomType
}

// Srid is the spatial reference identifier of the geom
func (e *Hdr) Srid() uint32 {
	return e.srid
}

// Dimension of the geometry
func (e *Hdr) Dimension() Dimension {
	return e.dim
}

// GeomType of the geom
func (e *Hdr) Type() GeomType {
	return e.gtype
}

// ----------------------------------------------------------------------------
// Point
// ----------------------------------------------------------------------------

// Point
type Point struct {
	Hdr
	Coordinate
}

func (p *Point) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"Point","coordinates":`)
	p.Coordinate.appendGeoJSON(&sb)
	sb.WriteString(`,`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)
	return sb.String()
}

func (p *Point) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";POINT")

	if p.dim == XYM {
		sb.WriteString("M")
	}
	sb.WriteString("(")
	sb.WriteString(strconv.FormatFloat(p.Coordinate[0], 'f', -1, 64))
	for idx := 1; idx < len(p.Coordinate); idx++ {
		sb.WriteString(" ")
		sb.WriteString(strconv.FormatFloat(p.Coordinate[idx], 'f', -1, 64))
	}
	sb.WriteString(")")
	return sb.String()
}

func (p *Point) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// MultiPoint
// ----------------------------------------------------------------------------

type MultiPoint struct {
	Hdr
	Points []Point
}

func (p *MultiPoint) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"MultiPoint","coordinates":[`)

	limit := len(p.Points) - 1
	for idx, point := range p.Points {
		point.Coordinate.appendGeoJSON(&sb)

		if idx < limit {
			sb.WriteString(",")
		}
	}

	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)

	return sb.String()
}

func (p *MultiPoint) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";MULTIPOINT()")

	return sb.String()
}

func (p *MultiPoint) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// LineString
// ----------------------------------------------------------------------------

type LineString struct {
	Hdr
	Coordinates []Coordinate
}

func (p *LineString) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"LineString", "coordinates":[`)
	limit := len(p.Coordinates) - 1
	for idx, coord := range p.Coordinates {
		coord.appendGeoJSON(&sb)

		if idx < limit {
			sb.WriteString(",")
		}
	}
	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)

	return sb.String()
}

func (p *LineString) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";LINESTRING()")

	return sb.String()
}

func (p *LineString) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// MultiLineString
// ----------------------------------------------------------------------------

type MultiLineString struct {
	Hdr
	LineStrings []LineString
}

func (p *MultiLineString) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"MultiLineString", "coordinates":[`)

	llimit := len(p.LineStrings) - 1
	for lidx, linestring := range p.LineStrings {
		sb.WriteString("[")

		limit := len(linestring.Coordinates) - 1
		for idx, coord := range linestring.Coordinates {
			coord.appendGeoJSON(&sb)

			if idx < limit {
				sb.WriteString(",")
			}
		}

		sb.WriteString("]")

		if lidx < llimit {
			sb.WriteString(",")
		}
	}

	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)

	return sb.String()
}

func (p *MultiLineString) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";MULTILINESTRING()")

	return sb.String()
}

func (p *MultiLineString) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// Polygon
// ----------------------------------------------------------------------------

type Polygon struct {
	Hdr
	Rings []LinearRing
}

func (p *Polygon) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"Polygon", "coordinates":[`)
	rlimit := len(p.Rings) - 1
	for ridx, lring := range p.Rings {
		lring.appendGeoJSON(&sb)

		if ridx < rlimit {
			sb.WriteString(",")
		}
	}
	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)
	return sb.String()
}

func (p *Polygon) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";POLYGON()")

	return sb.String()
}

func (p *Polygon) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// MultiPolygon
// ----------------------------------------------------------------------------

type MultiPolygon struct {
	Hdr
	Polygons []Polygon
}

func (p *MultiPolygon) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"MultiPolygon", "coordinates":[`)

	plimit := len(p.Polygons) - 1
	for pidx, polygon := range p.Polygons {
		sb.WriteString("[")

		rlimit := len(polygon.Rings) - 1
		for ridx, lring := range polygon.Rings {
			lring.appendGeoJSON(&sb)

			if ridx < rlimit {
				sb.WriteString(",")
			}
		}

		sb.WriteString("]")

		if pidx < plimit {
			sb.WriteString(",")
		}
	}

	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)
	return sb.String()
}

func (p *MultiPolygon) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";MULTIPOLYGON()")

	return sb.String()
}

func (p *MultiPolygon) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// GeometryCollection (a misnomer IMHO - should be called MultiGeometry)
// ----------------------------------------------------------------------------

type GeometryCollection struct {
	Hdr
	Geometries []Geometry
}

func (p *GeometryCollection) GeoJSON(crs, bbox bool) string {
	var sb bytes.Buffer
	sb.WriteString(`{"type":"GeometryCollection", "geometries":[`)
	limit := len(p.Geometries) - 1
	for idx, geom := range p.Geometries {
		sb.WriteString(geom.GeoJSON(false, false))

		if idx < limit {
			sb.WriteString(",")
		}
	}
	sb.WriteString(`],`)
	sb.WriteString(`"dim":"`)
	sb.WriteString(p.Hdr.dim.String())
	sb.WriteString(`"`)
	if crs {
		sb.WriteString(`,"crs":{"type":"name","properties":{"name":"EPSG:`)
		sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
		sb.WriteString(`"}}`)
	}
	sb.WriteString(`}`)

	return sb.String()
}

func (p *GeometryCollection) EWKT() string {
	var sb bytes.Buffer

	sb.WriteString("SRID=")
	sb.WriteString(strconv.FormatUint(uint64(p.srid), 10))
	sb.WriteString(";GEOMETRYCOLLECTION()")

	return sb.String()
}

func (p *GeometryCollection) String() string {
	return p.EWKT()
}

// ----------------------------------------------------------------------------
// LinearRing
// ----------------------------------------------------------------------------

type LinearRing struct {
	Coordinates []Coordinate
}

func (l *LinearRing) appendGeoJSON(sb *bytes.Buffer) {
	sb.WriteString("[")
	limit := len(l.Coordinates) - 1

	for idx, coord := range l.Coordinates {
		coord.appendGeoJSON(sb)

		if idx < limit {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
}

// ----------------------------------------------------------------------------
// Coordinate
// ----------------------------------------------------------------------------

// Coordinate type
type Coordinate []float64

func (c Coordinate) appendGeoJSON(sb *bytes.Buffer) {
	limit := len(c) - 1
	sb.WriteString("[")
	for idx, comp := range c {
		sb.WriteString(strconv.FormatFloat(comp, 'f', -1, 64))
		if idx < limit {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
}
