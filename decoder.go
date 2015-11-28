package ewkb

import (
	"encoding/binary"
	"fmt"
	"io"
)

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
	for idx = 0; idx < numGeoms; idx++ {
		d.u8()
		ghdr, err := unmarshalHdr(d)

		if err != nil {
			return nil, err
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
func unmarshal(d *decoder, hdr *Hdr) (Geometry, error) {
	switch hdr.gtype {
	case POINT:
		return unmarshalPoint(d, hdr)
	case LINESTRING:
		return unmarshalLineString(d, hdr)
	case POLYGON:
		return unmarshalPolygon(d, hdr)
	case MULTIPOINT:
		return unmarshalMultiPoint(d, hdr)
	case MULTILINESTRING:
		return unmarshalMultiLineString(d, hdr)
	case MULTIPOLYGON:
		return unmarshalMultiPolygon(d, hdr)
	case GEOMETRYCOLLECTION:
		return unmarshalGeometryCollection(d, hdr)
	default:
		return nil, ErrUnsupportedGeom
	}
}
