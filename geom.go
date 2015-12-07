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
