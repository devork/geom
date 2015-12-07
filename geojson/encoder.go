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

package geojson

import (
	"bytes"
	"io"
	"strconv"

	"github.com/devork/geom"
)

// snippets of geojson
var (
	pointHdr      = []byte(`{"type":"Point","coordinates":`)
	linestringHdr = []byte(`{"type":"LineString", "coordinates":`)
	polygonHdr    = []byte(`{"type":"Polygon", "coordinates":`)
	dquote        = []byte(`"`)
	comma         = []byte(`,`)
	lbrace        = []byte(`{`)
	rbrace        = []byte(`}`)
	lparen        = []byte(`[`)
	rparen        = []byte(`]`)
)

func Encode(g geom.Geometry, w io.Writer) error {
	switch g := g.(type) {
	case *geom.Point:
		return marshalPoint(g, w)
	case *geom.LineString:
		return marshalLineString(g, w)
	case *geom.Polygon:
		return marshalPolygon(g, w)
	default:
		return geom.ErrUnsupportedGeom

	}
}

func marshalPolygon(p *geom.Polygon, w io.Writer) error {
	var sb bytes.Buffer
	sb.Write(polygonHdr)
	sb.Write(lparen)
	rlimit := len(p.Rings) - 1
	for ridx, lring := range p.Rings {
		marshalLinearRing(&lring, &sb)

		if ridx < rlimit {
			sb.Write(comma)
		}
	}
	sb.Write(rparen)
	sb.Write(rbrace)

	_, err := w.Write(sb.Bytes())

	return err
}

func marshalLineString(ls *geom.LineString, w io.Writer) error {
	var sb bytes.Buffer
	sb.Write(linestringHdr)
	sb.Write(lparen)
	limit := len(ls.Coordinates) - 1
	for idx, coord := range ls.Coordinates {
		marshalCoord(&coord, &sb)

		if idx < limit {
			sb.Write(comma)
		}
	}
	sb.Write(rparen)
	sb.Write(rbrace)

	_, err := w.Write(sb.Bytes())

	return err
}

func marshalPoint(g *geom.Point, w io.Writer) error {

	var sb bytes.Buffer

	sb.Write(pointHdr)
	marshalCoord(&g.Coordinate, &sb)
	sb.Write(rbrace)
	_, err := w.Write(sb.Bytes())

	return err
}

func marshalLinearRing(l *geom.LinearRing, sb *bytes.Buffer) {

	sb.Write(lparen)
	limit := len(l.Coordinates) - 1

	for idx, coord := range l.Coordinates {
		marshalCoord(&coord, sb)

		if idx < limit {
			sb.Write(comma)
		}
	}
	sb.Write(rparen)

}

func marshalCoord(c *geom.Coordinate, sb *bytes.Buffer) {
	limit := len(*c) - 1
	sb.Write(lparen)
	for idx, comp := range *c {
		sb.Write([]byte(strconv.FormatFloat(comp, 'f', -1, 64)))
		if idx < limit {
			sb.Write(comma)
		}
	}
	sb.Write(rparen)
}
