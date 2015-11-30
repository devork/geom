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

import "github.com/devork/geom"

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
	xy      dimension = 0x0000
	xym     dimension = 0x4000
	xyz     dimension = 0x8000
	xyzm    dimension = 0xC000
	xys     dimension = 0x2000
	xyms    dimension = 0x6000
	xyzs    dimension = 0xA000
	xyzms   dimension = 0xE000
	unknown dimension = 0xFFFF
)

// Geometry types
const (
	geometry           geomtype = 0x0000
	point              geomtype = 0x0001
	linestring         geomtype = 0x0002
	polygon            geomtype = 0x0003
	multipoint         geomtype = 0x0004
	multilinestring    geomtype = 0x0005
	multipolygon       geomtype = 0x0006
	geometrycollection geomtype = 0x0007
	circularstring     geomtype = 0x0008
	compoundcurve      geomtype = 0x0009
	curvepolygon       geomtype = 0x000a
	multicurve         geomtype = 0x000b
	multisurface       geomtype = 0x000c
	curve              geomtype = 0x000d
	surface            geomtype = 0x000e
	polyhedralsurface  geomtype = 0x000f
	tin                geomtype = 0x0010
	triangle           geomtype = 0x0011
)

// ----------------------------------------------------------------------------
// Dimension
// ----------------------------------------------------------------------------

// Dimension of the geometry
type dimension uint16

func (d dimension) dim() geom.Dimension {
	switch d {
	case xy, xys:
		return geom.XY
	case xyz, xyzs:
		return geom.XYZ
	case xym, xyms:
		return geom.XYM
	case xyzm, xyzms:
		return geom.XYZM
	default:
		return geom.UNKNOWN
	}
}

// ----------------------------------------------------------------------------
// geomType
// ----------------------------------------------------------------------------

// geomType is the bitmask of the geom
type geomtype uint16

func (g geomtype) String() string {
	switch g {
	case geometry:
		return "GEOMETRY"
	case point:
		return "POINT"
	case linestring:
		return "LINESTRING"
	case polygon:
		return "POLYGON"
	case multipoint:
		return "MULTIPOINT"
	case multilinestring:
		return "MULTILINESTRING"
	case multipolygon:
		return "MULTIPOLYGON"
	case geometrycollection:
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
