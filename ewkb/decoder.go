package ewkb

import (
	"encoding/binary"
	"io"

	"github.com/devork/geom"
)

// wraps the byte order and reader into a single struct for easier mainpulation
type decoder struct {
	order  binary.ByteOrder
	reader io.Reader
	dim    dimension
	gtype  geomtype
	srid   uint32
}

func (d *decoder) read(data interface{}) error {
	return binary.Read(d.reader, d.order, data)
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

func (d *decoder) hdr() *geom.Hdr {
	return &geom.Hdr{Dim: d.dim.dim(), Srid: d.srid}
}

func newDecoder(r io.Reader) (*decoder, error) {
	var otype byte
	err := binary.Read(r, binary.BigEndian, &otype)

	if err != nil {
		return nil, err
	}

	if otype == bigEndian {
		return &decoder{order: binary.BigEndian, reader: r}, nil
	}
	return &decoder{order: binary.LittleEndian, reader: r}, nil

}

// Decode converts from the given reader to a geom type
func Decode(r io.Reader) (geom.Geometry, error) {

	decoder, err := newDecoder(r)

	if err != nil {
		return nil, err
	}

	err = unmarshalHdr(decoder)

	if err != nil {
		return nil, err
	}

	return unmarshal(decoder)
}

func unmarshalHdr(d *decoder) error {

	gtype, err := d.u32()

	if err != nil {
		return err
	}

	// fmt.Printf("gtype = %X\n", gtype)

	var geomv, dim uint16

	geomv = uint16(gtype & uint32(0xFFFF))
	dim = uint16(gtype >> 16)
	d.dim = dimension(dim)

	// fmt.Printf("geomv = %X [%[1]d]\n", geomv)
	// fmt.Printf("dim = %X\n", dim)

	// switch on the mask first to check for EWKB
	switch d.dim {
	case xys, xyms, xyzs, xyzms:
		err = d.read(&d.srid)

		if err != nil {
			return err
		}
	}

	// fallback check for WKB & reset the WKB versions to regular geometry types
	switch {
	case geomv >= wkbzm:
		geomv = geomv - wkbzm
		d.dim = xyzm
	case geomv >= wkbm:
		geomv = geomv - wkbm
		d.dim = xym
	case geomv >= wkbz:
		geomv = geomv - wkbz
		d.dim = xyz
	}

	d.gtype = geomtype(geomv)

	return nil
}

func unmarshalPoint(d *decoder) (geom.Geometry, error) {
	coord, err := unmarshalCoord(d)
	if err != nil {
		return nil, err
	}

	return &geom.Point{*d.hdr(), *coord}, nil
}

func unmarshalMultiPoint(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	points := make([]geom.Point, numPoints, numPoints)
	var point geom.Geometry
	var hdr = d.hdr()
	for idx = 0; idx < numPoints; idx++ {
		d.u8() // byteorder
		err = unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		point, err = unmarshalPoint(d)

		if err != nil {
			return nil, err
		}

		points[idx] = *point.(*geom.Point)
	}

	return &geom.MultiPoint{*hdr, points}, nil
}

func unmarshalLineString(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	coords := make([]geom.Coordinate, numPoints, numPoints)
	var coord *geom.Coordinate
	for idx = 0; idx < numPoints; idx++ {
		coord, err = unmarshalCoord(d)

		if err != nil {
			return nil, err
		}

		coords[idx] = *coord
	}

	return &geom.LineString{*d.hdr(), coords}, nil
}

func unmarshalMultiLineString(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numStrings, err := d.u32()

	if err != nil {
		return nil, err
	}

	lstrings := make([]geom.LineString, numStrings, numStrings)
	var lstring geom.Geometry
	var hdr = d.hdr()
	for idx = 0; idx < numStrings; idx++ {
		d.u8() // byteorder
		err = unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		lstring, err = unmarshalLineString(d)

		if err != nil {
			return nil, err
		}

		lstrings[idx] = *lstring.(*geom.LineString)
	}

	return &geom.MultiLineString{*hdr, lstrings}, nil
}

func unmarshalPolygon(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numRings, err := d.u32()

	if err != nil {
		return nil, err
	}

	rings := make([]geom.LinearRing, numRings, numRings)
	var ring *geom.LinearRing
	for idx = 0; idx < numRings; idx++ {
		ring, err = unmarshalLinearRing(d)

		if err != nil {
			return nil, err
		}

		rings[idx] = *ring
	}

	return &geom.Polygon{*d.hdr(), rings}, nil
}

func unmarshalMultiPolygon(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numPolys, err := d.u32()

	if err != nil {
		return nil, err
	}

	polys := make([]geom.Polygon, numPolys, numPolys)
	var poly geom.Geometry
	var hdr = d.hdr()
	for idx = 0; idx < numPolys; idx++ {
		d.u8() // byteorder
		err = unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		poly, err = unmarshalPolygon(d)

		if err != nil {
			return nil, err
		}

		polys[idx] = *poly.(*geom.Polygon)
	}

	return &geom.MultiPolygon{*hdr, polys}, nil
}

func unmarshalGeometryCollection(d *decoder) (geom.Geometry, error) {
	var idx uint32
	numGeoms, err := d.u32()

	if err != nil {
		return nil, err
	}

	geoms := make([]geom.Geometry, numGeoms, numGeoms)
	var g geom.Geometry
	var hdr = d.hdr()
	for idx = 0; idx < numGeoms; idx++ {
		d.u8()
		err := unmarshalHdr(d)

		if err != nil {
			return nil, err
		}

		g, err = unmarshal(d)

		if err != nil {
			return nil, err
		}
		geoms[idx] = g
	}

	return &geom.GeometryCollection{*hdr, geoms}, nil
}

func unmarshalLinearRing(d *decoder) (*geom.LinearRing, error) {
	var idx uint32
	numPoints, err := d.u32()

	if err != nil {
		return nil, err
	}

	coords := make([]geom.Coordinate, numPoints, numPoints)
	var coord *geom.Coordinate
	for idx = 0; idx < numPoints; idx++ {
		coord, err = unmarshalCoord(d)

		if err != nil {
			return nil, err
		}

		coords[idx] = *coord
	}

	return &geom.LinearRing{coords}, nil
}

func unmarshalCoord(d *decoder) (*geom.Coordinate, error) {
	var size int
	switch d.dim {
	case xyz, xyzs, xym, xyms:
		size = 3
	case xyzm, xyzms:
		size = 4
	default:
		size = 2
	}

	var coord geom.Coordinate = make([]float64, size, size)
	for idx := 0; idx < size; idx++ {
		err := d.read(&coord[idx])

		if err != nil {
			return nil, err
		}
	}

	return &coord, nil
}

// Resolves a geometry type to its unmarshaller instance
func unmarshal(d *decoder) (geom.Geometry, error) {
	switch d.gtype {
	case point:
		return unmarshalPoint(d)
	case linestring:
		return unmarshalLineString(d)
	case polygon:
		return unmarshalPolygon(d)
	case multipoint:
		return unmarshalMultiPoint(d)
	case multilinestring:
		return unmarshalMultiLineString(d)
	case multipolygon:
		return unmarshalMultiPolygon(d)
	case geometrycollection:
		return unmarshalGeometryCollection(d)
	default:
		return nil, ErrUnsupportedGeom
	}
}
