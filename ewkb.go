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
