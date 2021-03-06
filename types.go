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

// Geometry interface
type Geometry interface {
	Dimension() Dimension
	SRID() uint32
	Type() string
}

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

// Point
type Point struct {
	Hdr
	Coordinate
}

func (p *Point) Type() string {
	return "point"
}

// MultiPoint
type MultiPoint struct {
	Hdr
	Points []Point
}

func (m *MultiPoint) Type() string {
	return "multipoint"
}

// LineString
type LineString struct {
	Hdr
	Coordinates []Coordinate
}

func (l *LineString) Type() string {
	return "linestring"
}

// MultiLineString
type MultiLineString struct {
	Hdr
	LineStrings []LineString
}

func (l *MultiLineString) Type() string {
	return "multilinestring"
}

// Polygon
type Polygon struct {
	Hdr
	Rings []LinearRing
}

func (l *Polygon) Type() string {
	return "polygon"
}

// MultiPolygon
type MultiPolygon struct {
	Hdr
	Polygons []Polygon
}

func (l *MultiPolygon) Type() string {
	return "multipolygon"
}

// GeometryCollection (a misnomer IMHO - should be called MultiGeometry)
type GeometryCollection struct {
	Hdr
	Geometries []Geometry
}

func (l *GeometryCollection) Type() string {
	return "geometrycollection"
}

// LinearRing
type LinearRing struct {
	Coordinates []Coordinate
}

// Coordinate
type Coordinate []float64
