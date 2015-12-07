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
	default:
		return geom.ErrUnsupportedGeom

	}
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
