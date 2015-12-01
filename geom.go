package geom

import (
	"errors"
	"io"
)

// Common error types
var (
	ErrNoGeometry      = errors.New("no geometry specified")
	ErrUnsupportedGeom = errors.New("cannot encode unknown geometry")
	ErrUnknownDim      = errors.New("unknown dimension")
)

type Encoder interface {
	Encode(g *Geometry, w io.Writer) error
}

type Decoder interface {
	Decode(r io.Reader) (Geometry, error)
}
