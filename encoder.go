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

	err := marshalHdr(g, e)

	if err != nil {
		return err
	}

	switch g := g.(type) {
	case *Point:
		return marshalPoint(g, e)
	case *LineString:
		return marshalLineString(g, e)
	case *Polygon:
		return marshalPolygon(g, e)
	case *MultiPoint:
		return marshalMultiPoint(g, e)
	case *MultiLineString:
		return marshalMultiLineString(g, e)
	case *MultiPolygon:
		return marshalMultiPolygon(g, e)
	default:
		return ErrUnsupportedGeom
	}

	// case MULTILINESTRING:
	// 	return marshalMultiLineString
	// case MULTIPOLYGON:
	// 	return marshalMultiPolygon
	// case GEOMETRYCOLLECTION:
	// 	return marshalGeometryCollection
}

func marshalPoint(p *Point, e *encoder) error {
	return marshalCoord(&(*p).Coordinate, e)
}

func marshalMultiPoint(mp *MultiPoint, e *encoder) error {
	err := e.write(uint32(len(mp.Points)))

	if err != nil {
		return err
	}

	for _, point := range mp.Points {
		err = marshalHdr(&point, e)

		if err != nil {
			return err
		}

		err = marshalPoint(&point, e)

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalLineString(l *LineString, e *encoder) error {
	err := e.write(uint32(len(l.Coordinates)))

	if err != nil {
		return err
	}

	for _, coord := range l.Coordinates {
		err = marshalCoord(&coord, e)

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalMultiLineString(ml *MultiLineString, e *encoder) error {
	err := e.write(uint32(len(ml.LineStrings)))
	if err != nil {
		return err
	}

	for _, ls := range ml.LineStrings {
		err = marshalHdr(&ls, e)

		if err != nil {
			return err
		}

		err = marshalLineString(&ls, e)

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalPolygon(p *Polygon, e *encoder) error {
	err := e.write(uint32(len(p.Rings)))

	if err != nil {
		return err
	}

	for _, ring := range p.Rings {
		err = marshalLinearRing(&ring, e)

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalMultiPolygon(mp *MultiPolygon, e *encoder) error {
	err := e.write(uint32(len(mp.Polygons)))

	if err != nil {
		return err
	}

	for _, polygon := range mp.Polygons {

		err = marshalHdr(&polygon, e)

		if err != nil {
			return err
		}

		err = marshalPolygon(&polygon, e)

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalLinearRing(l *LinearRing, e *encoder) error {
	err := e.write(uint32(len(l.Coordinates)))

	if err != nil {
		return err
	}

	for _, coord := range l.Coordinates {
		err = marshalCoord(&coord, e)

		if err != nil {
			return err
		}
	}

	return nil
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

func marshalHdr(g Geometry, e *encoder) error {
	err := e.write(bigEndian)

	if err != nil {
		return err
	}

	var gtype uint32
	var writeSrid = false
	switch g.Dimension() {
	case XYS, XYZS, XYMS, XYZMS:
		gtype = uint32(g.Dimension())
		gtype <<= 16
		gtype |= uint32(g.Type())
		writeSrid = true
	case XYZM:
		gtype = uint32(wkbzm) + uint32(g.Type())
	case XYM:
		gtype = uint32(wkbm) + uint32(g.Type())
	case XYZ:
		gtype = uint32(wkbz) + uint32(g.Type())
	default:
		gtype = uint32(g.Type())
	}

	err = e.write(gtype)

	if err != nil {
		return err
	}

	if writeSrid {
		return e.write(g.Srid())
	}

	return nil
}
