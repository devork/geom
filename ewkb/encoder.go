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
	"io"

	"github.com/devork/geom"
)

type encoder struct {
	w io.Writer
	o binary.ByteOrder
}

func (d *encoder) write(data interface{}) error {
	return binary.Write(d.w, d.o, data)
}

//Encode will take the given geometry and write to the specified writer
func Encode(g geom.Geometry, w io.Writer) error {

	if g == nil {
		return geom.ErrNoGeometry
	}

	e := &encoder{w, binary.BigEndian}

	err := marshalHdr(g, e)

	if err != nil {
		return err
	}

	return marshal(g, e)
}

func marshalPoint(p *geom.Point, e *encoder) error {
	return marshalCoord(&(*p).Coordinate, e)
}

func marshalMultiPoint(mp *geom.MultiPoint, e *encoder) error {
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

func marshalLineString(l *geom.LineString, e *encoder) error {
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

func marshalMultiLineString(ml *geom.MultiLineString, e *encoder) error {
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

func marshalPolygon(p *geom.Polygon, e *encoder) error {
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

func marshalMultiPolygon(mp *geom.MultiPolygon, e *encoder) error {
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

func marshalGeometryCollection(gc *geom.GeometryCollection, e *encoder) error {
	err := e.write(uint32(len(gc.Geometries)))

	if err != nil {
		return err
	}

	for _, g := range gc.Geometries {
		err = marshalHdr(g, e)

		if err != nil {
			return err
		}

		err = marshal(g, e)

		if err != nil {
			return err
		}

	}

	return nil
}

func marshalLinearRing(l *geom.LinearRing, e *encoder) error {
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

func marshalCoord(c *geom.Coordinate, e *encoder) error {
	var err error
	for idx := 0; idx < len(*c); idx++ {
		err = e.write((*c)[idx])

		if err != nil {
			return err
		}
	}

	return nil
}

func marshalHdr(g geom.Geometry, e *encoder) error {
	err := e.write(bigEndian)

	if err != nil {
		return err
	}

	var field uint32
	var gtype geomtype
	var writeSrid bool

	switch g.(type) {
	case *geom.Point:
		gtype = point
	case *geom.LineString:
		gtype = linestring
	case *geom.Polygon:
		gtype = polygon
	case *geom.MultiPoint:
		gtype = multipoint
	case *geom.MultiLineString:
		gtype = multilinestring
	case *geom.MultiPolygon:
		gtype = multipolygon
	case *geom.GeometryCollection:
		gtype = geometrycollection
	default:
		return geom.ErrUnsupportedGeom
	}

	if g.SRID() != 0 {

		writeSrid = true

		switch g.Dimension() {
		case geom.XY:
			field = uint32(xys)
		case geom.XYZ:
			field = uint32(xyzs)
		case geom.XYM:
			field = uint32(xyms)
		case geom.XYZM:
			field = uint32(xyzms)
		default:
			return geom.ErrUnknownDim
		}

		field <<= 16
		field |= uint32(gtype)
	} else {
		switch g.Dimension() {
		case geom.XYZM:
			field = uint32(wkbzm) + uint32(gtype)
		case geom.XYM:
			field = uint32(wkbm) + uint32(gtype)
		case geom.XYZ:
			field = uint32(wkbz) + uint32(gtype)
		case geom.XY:
			field = uint32(gtype)
		default:
			return geom.ErrUnknownDim
		}
	}

	err = e.write(field)

	if err != nil {
		return err
	}

	if writeSrid {
		return e.write(g.SRID())
	}

	return nil
}

func marshal(g geom.Geometry, e *encoder) error {
	switch g := g.(type) {
	case *geom.Point:
		return marshalPoint(g, e)
	case *geom.LineString:
		return marshalLineString(g, e)
	case *geom.Polygon:
		return marshalPolygon(g, e)
	case *geom.MultiPoint:
		return marshalMultiPoint(g, e)
	case *geom.MultiLineString:
		return marshalMultiLineString(g, e)
	case *geom.MultiPolygon:
		return marshalMultiPolygon(g, e)
	case *geom.GeometryCollection:
		return marshalGeometryCollection(g, e)
	default:
		return geom.ErrUnsupportedGeom
	}
}
