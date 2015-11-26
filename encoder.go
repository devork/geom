package ewkb

import (
	"encoding/binary"
	"errors"
	"io"
)

type encoder struct {
	w io.Writer
	o binary.ByteOrder
}

func (d *encoder) write(data interface{}) error {
	return binary.Write(d.w, d.o, data)
}

// Common error types
var (
	ErrNoGeometry      = errors.New("no geometry specified")
	ErrUnsupportedGeom = errors.New("cannot encode unknown geometry ")
)

//Encode will take the given geometry and write to the specified writer
func Encode(g Geometry, w io.Writer) error {

	if g == nil {
		return ErrNoGeometry
	}

	e := &encoder{w, binary.BigEndian}

	switch g.Type() {
	case POINT:
		return marshalPoint(g.(*Point), e)
	// case LINESTRING:
	// 	return marshalLineString
	// case POLYGON:
	// 	return marshalPolygon
	// case MULTIPOINT:
	// 	return marshalMultiPoint
	// case MULTILINESTRING:
	// 	return marshalMultiLineString
	// case MULTIPOLYGON:
	// 	return marshalMultiPolygon
	// case GEOMETRYCOLLECTION:
	// 	return marshalGeometryCollection
	default:
		return ErrUnsupportedGeom
	}

}

func marshalPoint(p *Point, e *encoder) error {
	err := writeHeader(p, e)

	if err != nil {
		return err
	}

	return marshalCoord(&(*p).Coordinate, e)
}

func marshalCoord(c *Coordinate, e *encoder) error {
	var err error
	for idx := 0; idx < len(*c); idx++ {
		err = e.write((*c)[idx])

		if err != nil {
			return err
		}
	}

	return nil
}

func writeHeader(g Geometry, e *encoder) error {
	err := e.write(bigEndian)

	if err != nil {
		return err
	}

	var gtype uint32

	if g.Srid() != 0 {
		// EWKB
		gtype = uint32(g.Dimension())
		gtype <<= 16
		gtype |= uint32(g.Type())
	} else {
		// WKB
		switch g.Dimension() {
		case XYZM:
			gtype += uint32(wkbzm)
		case XYM:
			gtype += uint32(wkbm)
		case XYZ:
			gtype += uint32(wkbz)
		}

		gtype += uint32(g.Type())
	}

	err = e.write(gtype)

	if err != nil {
		return err
	}

	if g.Srid() != 0 {
		return e.write(g.Srid())
	}

	return nil
}
