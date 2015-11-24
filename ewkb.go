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
	"encoding/binary"
	"fmt"
	"io"
)

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

// https://trac.osgeo.org/postgis/browser/trunk/doc/ZMSgeoms.txt
//
// Supported dimensions:
//
//		XY		2 dimensional
//		XYZ	 	3 dimensional
//		XYZM	4 dimensional
//
// 2.1. Definition of ZM-Geometry
//
//	a) A geometry can have either 2, 3 or 4 dimensions.
//	b) 3rd dimension of a 3d geometry can either represent Z or M (3DZ or 3DM).
//	c) 4d geometries contain both Z and M (in this order).
//	d) M and Z values are associated with every vertex.
//	e) M and Z values are undefined within surface interiors.
//
//	Any ZM-Geometry can be converted into a 2D geometry by discarding all its
//	Z and M values. The resulting 2D geometry is the "shadow" of the ZM-Geometry.
//	2D geometries cannot be safely converted into ZM-Geometries, since their Z
//	and M values are undefined, and not necessarily zero.
//
// These constants also represent the bit masks for the EWKB dims
const (
	XY    Dimension = 0x0000
	XYM   Dimension = 0x4000
	XYZ   Dimension = 0x8000
	XYZM  Dimension = 0xC000
	XYS   Dimension = 0x2000
	XYMS  Dimension = 0x6000
	XYZS  Dimension = 0xA000
	XYZMS Dimension = 0xE000
)

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
	case CIRCULARSTRING:
		return "CIRCULARSTRING"
	case COMPOUNDCURVE:
		return "COMPOUNDCURVE"
	case CURVEPOLYGON:
		return "CURVEPOLYGON"
	case MULTICURVE:
		return "MULTICURVE"
	case MULTISURFACE:
		return "MULTISURFACE"
	case CURVE:
		return "CURVE"
	case SURFACE:
		return "SURFACE"
	case POLYHEDRALSURFACE:
		return "POLYHEDRALSURFACE "
	case TIN:
		return "TIN"
	case TRIANGLE:
		return "TRIANGLE"
	default:
		return "UNKNOWN"
	}
}

// Geometry types
const (
	GEOMETRY           GeomType = 0x0000
	POINT              GeomType = 0x0001
	LINESTRING         GeomType = 0x0002
	POLYGON            GeomType = 0x0003
	MULTIPOINT         GeomType = 0x0004
	MULTILINESTRING    GeomType = 0x0005
	MULTIPOLYGON       GeomType = 0x0006
	GEOMETRYCOLLECTION GeomType = 0x0007
	CIRCULARSTRING     GeomType = 0x0008
	COMPOUNDCURVE      GeomType = 0x0009
	CURVEPOLYGON       GeomType = 0x000A
	MULTICURVE         GeomType = 0x000B
	MULTISURFACE       GeomType = 0x000C
	CURVE              GeomType = 0x000D
	SURFACE            GeomType = 0x000E
	POLYHEDRALSURFACE  GeomType = 0x000F
	TIN                GeomType = 0x0010
	TRIANGLE           GeomType = 0x0011
)

// WKB extensions for Z, M, and ZM. these extensions are applied to the base geometry types,
// such that a ZM version of a Point = 17 + 3000 = 3017 (0xBC9)
const (
	wkbz  uint16 = 1000
	wkbm  uint16 = 2000
	wkbzm uint16 = 3000
)

// Big or Little endian identifiers
const (
	bigEndian uint8 = 0x00
)

// handler type to convert a stream of bytes to a geometry type
type unmarshaller func(d *decoder, h *Hdr) (Geometry, error)

// wraps the byte order and reader into a single struct for easier mainpulation
type decoder struct {
	o binary.ByteOrder
	r io.Reader
}

func (d *decoder) read(data interface{}) error {
	return binary.Read(d.r, d.o, data)
}

func (d *decoder) u8() (uint8, error) {
	var value uint8

	err := d.read(&value)

	return value, err
}

func (d *decoder) u16() (uint16, error) {
	var value uint16

	err := d.read(&value)

	return value, err
}

func (d *decoder) u32() (uint32, error) {
	var value uint32

	err := d.read(&value)

	return value, err
}

func newDecoder(r io.Reader) (*decoder, error) {
	var otype byte
	err := binary.Read(r, binary.BigEndian, &otype)

	if err != nil {
		return nil, err
	}

	if otype == bigEndian {
		return &decoder{binary.BigEndian, r}, nil
	}
	return &decoder{binary.LittleEndian, r}, nil

}

// Decode converts from the given reader to a geom type
func Decode(r io.Reader) (Geometry, error) {

	decoder, err := newDecoder(r)

	if err != nil {
		return nil, err
	}

	hdr, err := unmarshalHdr(decoder)

	if err != nil {
		return nil, err
	}

	unmarshal := resolve(hdr.gtype)

	if unmarshal == nil {
		return nil, fmt.Errorf("Unknown geometry type: %X", hdr.gtype)
	}

	return unmarshal(decoder, hdr)
}

func unmarshalHdr(d *decoder) (*Hdr, error) {

	hdr := Hdr{}

	gtype, err := d.u32()

	if err != nil {
		return nil, err
	}

	var geom, mask uint16

	geom = uint16(gtype & uint32(0xFFFF))
	mask = uint16(gtype >> 16)

	// switch on the mask first to check for EWKB
	switch Dimension(mask) {
	case XYS, XYMS, XYZS, XYZMS:
		err = d.read(&hdr.srid)

		if err != nil {
			return nil, err
		}
		hdr.dim = Dimension(mask)
	case XY, XYM, XYZ, XYZM:
		hdr.dim = Dimension(mask)
	default:
		return nil, fmt.Errorf("Unknown EWKB dimension type: %X", mask)
	}

	// fallback check for WKB & reset the WKB versions to regular geometry types
	switch {
	case geom >= wkbzm:
		hdr.dim = XYZM
		geom = geom - wkbzm
	case geom >= wkbm:
		hdr.dim = XYM
		geom = geom - wkbm
	case geom >= wkbz:
		hdr.dim = XYZ
		geom = geom - wkbz
	}

	hdr.gtype = GeomType(geom)

	return &hdr, nil
}

func unmarshalPoint(d *decoder, hdr *Hdr) (Geometry, error) {
	coord, err := unmarshalCoord(d, hdr.dim)
	if err != nil {
		return nil, err
	}

	return &Point{*hdr, *coord}, nil
}

func unmarshalMultiPoint(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	points := make([]Point, numPoints, numPoints)
	var point Geometry

	for idx = 0; idx < numPoints; idx++ {
		d.u8() // byteorder
		phdr, err := unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		point, err = unmarshalPoint(d, phdr)

		if err != nil {
			return nil, err
		}

		points[idx] = *point.(*Point)
	}

	return &MultiPoint{*hdr, points}, nil
}

func unmarshalLineString(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	coords := make([]Coordinate, numPoints, numPoints)
	var coord *Coordinate
	for idx = 0; idx < numPoints; idx++ {
		coord, err = unmarshalCoord(d, hdr.dim)

		if err != nil {
			return nil, err
		}

		coords[idx] = *coord
	}

	return &LineString{*hdr, coords}, nil
}

func unmarshalMultiLineString(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numStrings, err := d.u32()

	if err != nil {
		return nil, err
	}

	lstrings := make([]LineString, numStrings, numStrings)
	var lstring Geometry
	for idx = 0; idx < numStrings; idx++ {
		d.u8() // byteorder
		lshdr, err := unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		lstring, err = unmarshalLineString(d, lshdr)

		if err != nil {
			return nil, err
		}

		lstrings[idx] = *lstring.(*LineString)
	}

	return &MultiLineString{*hdr, lstrings}, nil
}

func unmarshalPolygon(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numRings, err := d.u32()

	if err != nil {
		return nil, err
	}

	rings := make([]LinearRing, numRings, numRings)
	var ring *LinearRing
	for idx = 0; idx < numRings; idx++ {
		ring, err = unmarshalLinearRing(d, hdr)

		if err != nil {
			return nil, err
		}

		rings[idx] = *ring
	}

	return &Polygon{*hdr, rings}, nil
}

func unmarshalMultiPolygon(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numPolys, err := d.u32()

	if err != nil {
		return nil, err
	}

	polys := make([]Polygon, numPolys, numPolys)
	var poly Geometry
	for idx = 0; idx < numPolys; idx++ {
		d.u8() // byteorder
		phdr, err := unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		poly, err = unmarshalPolygon(d, phdr)

		if err != nil {
			return nil, err
		}

		polys[idx] = *poly.(*Polygon)
	}

	return &MultiPolygon{*hdr, polys}, nil
}

func unmarshalGeometryCollection(d *decoder, hdr *Hdr) (Geometry, error) {
	var idx uint32
	numGeoms, err := d.u32()

	if err != nil {
		return nil, err
	}

	geoms := make([]Geometry, numGeoms, numGeoms)
	var geom Geometry
	var unmarshal unmarshaller
	for idx = 0; idx < numGeoms; idx++ {
		d.u8()
		ghdr, err := unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		unmarshal = resolve(ghdr.gtype)

		if unmarshal == nil {
			return nil, fmt.Errorf("Unknown geometry type: %X", ghdr.gtype)
		}

		geom, err = unmarshal(d, ghdr)

		if err != nil {
			return nil, err
		}
		geoms[idx] = geom
	}

	return &GeometryCollection{*hdr, geoms}, nil
}

func unmarshalLinearRing(d *decoder, hdr *Hdr) (*LinearRing, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	coords := make([]Coordinate, numPoints, numPoints)
	var coord *Coordinate
	for idx = 0; idx < numPoints; idx++ {
		coord, err = unmarshalCoord(d, hdr.dim)

		if err != nil {
			return nil, err
		}

		coords[idx] = *coord
	}

	return &LinearRing{coords}, nil
}

func unmarshalCoord(d *decoder, dim Dimension) (*Coordinate, error) {
	var size int
	switch dim {
	case XYZ, XYZS, XYM, XYMS:
		size = 3
	case XYZM, XYZMS:
		size = 4
	default:
		size = 2
	}

	var coord Coordinate = make([]float64, size, size)
	for idx := 0; idx < size; idx++ {
		err := d.read(&coord[idx])

		if err != nil {
			return nil, err
		}
	}

	return &coord, nil
}

// Resolves a geometry type to its unmarshaller instance
func resolve(gtype GeomType) unmarshaller {
	switch gtype {
	case POINT:
		return unmarshalPoint
	case LINESTRING:
		return unmarshalLineString
	case POLYGON:
		return unmarshalPolygon
	case MULTIPOINT:
		return unmarshalMultiPoint
	case MULTILINESTRING:
		return unmarshalMultiLineString
	case MULTIPOLYGON:
		return unmarshalMultiPolygon
	case GEOMETRYCOLLECTION:
		return unmarshalGeometryCollection
	default:
		return nil
	}
}
