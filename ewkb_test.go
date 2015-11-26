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

/*
# Test Case Generation

All test cases were generated from PostGIS:

	LINESTRING (30 10, 10 30, 40 40)

	with g as (
		select st_geomfromewkt('SRID=27700;LINESTRING (30 10, 10 30, 40 40)') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))

	with g as (
		select st_geomfromewkt('SRID=27700;POLYGON ((30 10, 40 40, 20 40, 10 20, 30 10))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	POLYGON ((35 10, 45 45, 15 40, 10 20, 35 10),(20 30, 35 35, 30 20, 20 30))

	with g as (
		select st_geomfromewkt('SRID=27700;POLYGON ((35 10, 45 45, 15 40, 10 20, 35 10),(20 30, 35 35, 30 20, 20 30))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	MULTIPOINT ((10 40), (40 30), (20 20), (30 10))

	with g as (
		select st_geomfromewkt('SRID=27700;MULTIPOINT ((10 40), (40 30), (20 20), (30 10))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	MULTILINESTRING ((10 10, 20 20, 10 40), (40 40, 30 30, 40 20, 30 10))

	with g as (
		select st_geomfromewkt('SRID=27700;MULTILINESTRING ((10 10, 20 20, 10 40), (40 40, 30 30, 40 20, 30 10))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	MULTIPOLYGON (((40 40, 20 45, 45 30, 40 40)),((20 35, 10 30, 10 10, 30 5, 45 20, 20 35),(30 20, 20 15, 20 25, 30 20)))

	with g as (
		select st_geomfromewkt('SRID=27700;MULTIPOLYGON (((40 40, 20 45, 45 30, 40 40)),((20 35, 10 30, 10 10, 30 5, 45 20, 20 35),(30 20, 20 15, 20 25, 30 20)))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

	GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))

	with g as (
		select st_geomfromewkt('SRID=27700;GEOMETRYCOLLECTION(POINT(4 6),LINESTRING(4 6,7 10))') as geom
	)
	select
		ST_AsEWKT(g.geom) as ewkt,
		st_asewkb(g.geom, 'XDR') as ewkb_xdr,
		st_asewkb(g.geom, 'HDR') as ewkb_hdr,
		st_asbinary(g.geom, 'XDR') as wkb_xdr,
		st_asbinary(g.geom, 'HDR') as wkb_hdr
	from
		g
	;

*/

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiGeometry(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000700006c340000000200000000014010000000000000401800000000000000000000020000000240100000000000004018000000000000401c0000000000004024000000000000", 27700, XYS, GEOMETRYCOLLECTION}, // ewkb_xdr
		{"0107000020346c000002000000010100000000000000000010400000000000001840010200000002000000000000000000104000000000000018400000000000001c400000000000002440", 27700, XYS, GEOMETRYCOLLECTION}, // ewkb_hdr
		{"00000000070000000200000000014010000000000000401800000000000000000000020000000240100000000000004018000000000000401c0000000000004024000000000000", 0, XY, GEOMETRYCOLLECTION},              // wkb_hdr
		{"010700000002000000010100000000000000000010400000000000001840010200000002000000000000000000104000000000000018400000000000001c400000000000002440", 0, XY, GEOMETRYCOLLECTION},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatalf("Failed to convert geom: error = %s", err)
		}

		gcol := geom.(*GeometryCollection)

		t.Log(gcol)

		assert.Equal(t, dataset.dim, gcol.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, gcol.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, gcol.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 2, len(gcol.Geometries))

		point := gcol.Geometries[0].(*Point)

		assert.InDelta(t, 4, point.Coordinate[0], 1e-9)
		assert.InDelta(t, 6, point.Coordinate[1], 1e-9)

		lstring := gcol.Geometries[1].(*LineString)

		assert.Equal(t, 2, len(lstring.Coordinates))

		assert.InDelta(t, 4, lstring.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 6, lstring.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 7, lstring.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 10, lstring.Coordinates[1][1], 1e-9)
	}
}

func TestMultiPolygon(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000600006c34000000020000000003000000010000000440440000000000004044000000000000403400000000000040468000000000004046800000000000403e0000000000004044000000000000404400000000000000000000030000000200000006403400000000000040418000000000004024000000000000403e00000000000040240000000000004024000000000000403e0000000000004014000000000000404680000000000040340000000000004034000000000000404180000000000000000004403e00000000000040340000000000004034000000000000402e00000000000040340000000000004039000000000000403e0000000000004034000000000000", 27700, XYS, MULTIPOLYGON}, // ewkb_xdr
		{"0106000020346c00000200000001030000000100000004000000000000000000444000000000000044400000000000003440000000000080464000000000008046400000000000003e4000000000000044400000000000004440010300000002000000060000000000000000003440000000000080414000000000000024400000000000003e40000000000000244000000000000024400000000000003e4000000000000014400000000000804640000000000000344000000000000034400000000000804140040000000000000000003e40000000000000344000000000000034400000000000002e40000000000000344000000000000039400000000000003e400000000000003440", 27700, XYS, MULTIPOLYGON}, // ewkb_hdr
		{"0000000006000000020000000003000000010000000440440000000000004044000000000000403400000000000040468000000000004046800000000000403e0000000000004044000000000000404400000000000000000000030000000200000006403400000000000040418000000000004024000000000000403e00000000000040240000000000004024000000000000403e0000000000004014000000000000404680000000000040340000000000004034000000000000404180000000000000000004403e00000000000040340000000000004034000000000000402e00000000000040340000000000004039000000000000403e0000000000004034000000000000", 0, XY, MULTIPOLYGON},              // wkb_hdr
		{"01060000000200000001030000000100000004000000000000000000444000000000000044400000000000003440000000000080464000000000008046400000000000003e4000000000000044400000000000004440010300000002000000060000000000000000003440000000000080414000000000000024400000000000003e40000000000000244000000000000024400000000000003e4000000000000014400000000000804640000000000000344000000000000034400000000000804140040000000000000000003e40000000000000344000000000000034400000000000002e40000000000000344000000000000039400000000000003e400000000000003440", 0, XY, MULTIPOLYGON},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse MultiPolygon: err = ", err)
		}

		mpolygon := geom.(*MultiPolygon)

		t.Log(mpolygon)

		assert.Equal(t, dataset.dim, mpolygon.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, mpolygon.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, mpolygon.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 2, len(mpolygon.Polygons))

		polygon := mpolygon.Polygons[0]
		assert.Equal(t, 1, len(polygon.Rings))

		lr := polygon.Rings[0]

		assert.Equal(t, 4, len(lr.Coordinates))

		assert.InDelta(t, 40, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 40, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 45, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 45, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 30, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 40, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 40, lr.Coordinates[3][1], 1e-9)

		polygon = mpolygon.Polygons[1]
		assert.Equal(t, 2, len(polygon.Rings))

		lr = polygon.Rings[0]

		assert.Equal(t, 6, len(lr.Coordinates))

		assert.InDelta(t, 20, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 35, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 10, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 30, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 10, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 10, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 30, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 5, lr.Coordinates[3][1], 1e-9)

		assert.InDelta(t, 45, lr.Coordinates[4][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[4][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[5][0], 1e-9)
		assert.InDelta(t, 35, lr.Coordinates[5][1], 1e-9)

		lr = polygon.Rings[1]

		assert.Equal(t, 4, len(lr.Coordinates))

		assert.InDelta(t, 30, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 15, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 25, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 30, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[3][1], 1e-9)

	}
}

func TestMultiLineString(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000500006c340000000200000000020000000340240000000000004024000000000000403400000000000040340000000000004024000000000000404400000000000000000000020000000440440000000000004044000000000000403e000000000000403e00000000000040440000000000004034000000000000403e0000000000004024000000000000", 27700, XYS, MULTILINESTRING}, // ewkb_xdr
		{"0105000020346c000002000000010200000003000000000000000000244000000000000024400000000000003440000000000000344000000000000024400000000000004440010200000004000000000000000000444000000000000044400000000000003e400000000000003e40000000000000444000000000000034400000000000003e400000000000002440", 27700, XYS, MULTILINESTRING}, // ewkb_hdr
		{"00000000050000000200000000020000000340240000000000004024000000000000403400000000000040340000000000004024000000000000404400000000000000000000020000000440440000000000004044000000000000403e000000000000403e00000000000040440000000000004034000000000000403e0000000000004024000000000000", 0, XY, MULTILINESTRING},              // wkb_hdr
		{"010500000002000000010200000003000000000000000000244000000000000024400000000000003440000000000000344000000000000024400000000000004440010200000004000000000000000000444000000000000044400000000000003e400000000000003e40000000000000444000000000000034400000000000003e400000000000002440", 0, XY, MULTILINESTRING},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse MultiLineString: err = ", err)
		}

		mlstring := geom.(*MultiLineString)

		t.Log(mlstring)

		assert.Equal(t, dataset.dim, mlstring.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, mlstring.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, mlstring.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 2, len(mlstring.LineStrings))

		lstring := mlstring.LineStrings[0]

		assert.Equal(t, 3, len(lstring.Coordinates))

		assert.InDelta(t, 10, lstring.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 10, lstring.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 20, lstring.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 20, lstring.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 10, lstring.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 40, lstring.Coordinates[2][1], 1e-9)

		lstring = mlstring.LineStrings[1]

		assert.Equal(t, 4, len(lstring.Coordinates))

		assert.InDelta(t, 40, lstring.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 40, lstring.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 30, lstring.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 30, lstring.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 40, lstring.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 20, lstring.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 30, lstring.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 10, lstring.Coordinates[3][1], 1e-9)
	}
}

func TestMultiPoint(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000400006c340000000400000000014024000000000000404400000000000000000000014044000000000000403e0000000000000000000001403400000000000040340000000000000000000001403e0000000000004024000000000000", 27700, XYS, MULTIPOINT}, // ewkb_xdr
		{"0104000020346c000004000000010100000000000000000024400000000000004440010100000000000000000044400000000000003e4001010000000000000000003440000000000000344001010000000000000000003e400000000000002440", 27700, XYS, MULTIPOINT}, // ewkb_hdr
		{"00000000040000000400000000014024000000000000404400000000000000000000014044000000000000403e0000000000000000000001403400000000000040340000000000000000000001403e0000000000004024000000000000", 0, XY, MULTIPOINT},              // wkb_hdr
		{"010400000004000000010100000000000000000024400000000000004440010100000000000000000044400000000000003e4001010000000000000000003440000000000000344001010000000000000000003e400000000000002440", 0, XY, MULTIPOINT},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse MultiPoint: err = ", err)
		}

		mpoint := geom.(*MultiPoint)

		t.Log(mpoint)

		assert.Equal(t, dataset.dim, mpoint.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, mpoint.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, mpoint.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 4, len(mpoint.Points))

		assert.InDelta(t, 10, mpoint.Points[0].Coordinate[0], 1e-9)
		assert.InDelta(t, 40, mpoint.Points[0].Coordinate[1], 1e-9)

		assert.InDelta(t, 40, mpoint.Points[1].Coordinate[0], 1e-9)
		assert.InDelta(t, 30, mpoint.Points[1].Coordinate[1], 1e-9)

		assert.InDelta(t, 20, mpoint.Points[2].Coordinate[0], 1e-9)
		assert.InDelta(t, 20, mpoint.Points[2].Coordinate[1], 1e-9)

		assert.InDelta(t, 30, mpoint.Points[3].Coordinate[0], 1e-9)
		assert.InDelta(t, 10, mpoint.Points[3].Coordinate[1], 1e-9)
	}
}

func TestHolePolygon(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000300006c3400000002000000054041800000000000402400000000000040468000000000004046800000000000402e00000000000040440000000000004024000000000000403400000000000040418000000000004024000000000000000000044034000000000000403e00000000000040418000000000004041800000000000403e00000000000040340000000000004034000000000000403e000000000000", 27700, XYS, POLYGON}, // ewkb_xdr
		{"0103000020346c0000020000000500000000000000008041400000000000002440000000000080464000000000008046400000000000002e40000000000000444000000000000024400000000000003440000000000080414000000000000024400400000000000000000034400000000000003e40000000000080414000000000008041400000000000003e40000000000000344000000000000034400000000000003e40", 27700, XYS, POLYGON}, // ewkb_hdr
		{"000000000300000002000000054041800000000000402400000000000040468000000000004046800000000000402e00000000000040440000000000004024000000000000403400000000000040418000000000004024000000000000000000044034000000000000403e00000000000040418000000000004041800000000000403e00000000000040340000000000004034000000000000403e000000000000", 0, XY, POLYGON},              // wkb_hdr
		{"0103000000020000000500000000000000008041400000000000002440000000000080464000000000008046400000000000002e40000000000000444000000000000024400000000000003440000000000080414000000000000024400400000000000000000034400000000000003e40000000000080414000000000008041400000000000003e40000000000000344000000000000034400000000000003e40", 0, XY, POLYGON},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse HolePolygon: err = ", err)
		}

		polygon := geom.(*Polygon)

		t.Log(polygon)

		assert.Equal(t, dataset.dim, polygon.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, polygon.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, polygon.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 2, len(polygon.Rings))

		lr := polygon.Rings[0]

		assert.Equal(t, 5, len(lr.Coordinates))

		assert.InDelta(t, 35, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 10, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 45, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 45, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 15, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 40, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 10, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[3][1], 1e-9)

		assert.InDelta(t, 35, lr.Coordinates[4][0], 1e-9)
		assert.InDelta(t, 10, lr.Coordinates[4][1], 1e-9)

		lr = polygon.Rings[1]

		assert.Equal(t, 4, len(lr.Coordinates))

		assert.InDelta(t, 20, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 30, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 35, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 35, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 30, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 30, lr.Coordinates[3][1], 1e-9)
	}
}

func TestPolygon(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000300006c340000000100000005403e0000000000004024000000000000404400000000000040440000000000004034000000000000404400000000000040240000000000004034000000000000403e0000000000004024000000000000", 27700, XYS, POLYGON}, // ewkb_xdr
		{"0103000020346c000001000000050000000000000000003e4000000000000024400000000000004440000000000000444000000000000034400000000000004440000000000000244000000000000034400000000000003e400000000000002440", 27700, XYS, POLYGON}, // ewkb_hdr
		{"00000000030000000100000005403e0000000000004024000000000000404400000000000040440000000000004034000000000000404400000000000040240000000000004034000000000000403e0000000000004024000000000000", 0, XY, POLYGON},              // wkb_hdr
		{"010300000001000000050000000000000000003e4000000000000024400000000000004440000000000000444000000000000034400000000000004440000000000000244000000000000034400000000000003e400000000000002440", 0, XY, POLYGON},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse Polygon: err = ", err)
		}

		polygon := geom.(*Polygon)

		t.Log(polygon)

		assert.Equal(t, dataset.dim, polygon.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, polygon.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, polygon.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)

		assert.Equal(t, 1, len(polygon.Rings))

		lr := polygon.Rings[0]

		assert.Equal(t, 5, len(lr.Coordinates))

		assert.InDelta(t, 30, lr.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 10, lr.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 40, lr.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 40, lr.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 20, lr.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 40, lr.Coordinates[2][1], 1e-9)

		assert.InDelta(t, 10, lr.Coordinates[3][0], 1e-9)
		assert.InDelta(t, 20, lr.Coordinates[3][1], 1e-9)

		assert.InDelta(t, 30, lr.Coordinates[4][0], 1e-9)
		assert.InDelta(t, 10, lr.Coordinates[4][1], 1e-9)
	}
}

func TestLineString(t *testing.T) {
	datasets := []struct {
		data  string
		srid  uint32
		dim   Dimension
		gtype GeomType
	}{
		{"002000000200006c3400000003403e00000000000040240000000000004024000000000000403e00000000000040440000000000004044000000000000", 27700, XYS, LINESTRING}, // ewkb_xdr
		{"0102000020346c0000030000000000000000003e40000000000000244000000000000024400000000000003e4000000000000044400000000000004440", 27700, XYS, LINESTRING}, // ewkb_hdr
		{"000000000200000003403e00000000000040240000000000004024000000000000403e00000000000040440000000000004044000000000000", 0, XY, LINESTRING},              // wkb_hdr
		{"0102000000030000000000000000003e40000000000000244000000000000024400000000000003e4000000000000044400000000000004440", 0, XY, LINESTRING},              // wkb_xdr
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse LineString: err = ", err)
		}

		lstring := geom.(*LineString)

		t.Log(lstring)

		assert.Equal(t, dataset.dim, lstring.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, lstring.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, lstring.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)
		assert.Equal(t, 3, len(lstring.Coordinates))

		assert.InDelta(t, 30, lstring.Coordinates[0][0], 1e-9)
		assert.InDelta(t, 10, lstring.Coordinates[0][1], 1e-9)

		assert.InDelta(t, 10, lstring.Coordinates[1][0], 1e-9)
		assert.InDelta(t, 30, lstring.Coordinates[1][1], 1e-9)

		assert.InDelta(t, 40, lstring.Coordinates[2][0], 1e-9)
		assert.InDelta(t, 40, lstring.Coordinates[2][1], 1e-9)
	}
}

func TestDimensionsAndEndian(t *testing.T) {

	datasets := []struct {
		data   string
		srid   uint32
		dim    Dimension
		gtype  GeomType
		points []float64
	}{
		// XY
		{"00000000013ff00000000000003ff0000000000000", 0, XY, POINT, []float64{1, 1}}, // EWKB XDR
		{"0101000000000000000000f03f000000000000f03f", 0, XY, POINT, []float64{1, 1}}, // EWKB HDR
		{"00000000013ff00000000000003ff0000000000000", 0, XY, POINT, []float64{1, 1}}, // WKB XDR
		{"0101000000000000000000f03f000000000000f03f", 0, XY, POINT, []float64{1, 1}}, // WKB HDR

		// XYS
		{"002000000100006c343ff00000000000003ff0000000000000", 27700, XYS, POINT, []float64{1, 1}}, // EWKB XDR
		{"0101000020346c0000000000000000f03f000000000000f03f", 27700, XYS, POINT, []float64{1, 1}}, // EWKB HDR
		{"00000000013ff00000000000003ff0000000000000", 0, XY, POINT, []float64{1, 1}},              // WKB XDR
		{"0101000000000000000000f03f000000000000f03f", 0, XY, POINT, []float64{1, 1}},              // WKB HDR

		// XYZ
		{"00800000013ff00000000000003ff00000000000003ff0000000000000", 0, XYZ, POINT, []float64{1, 1, 1}}, // EWKB XDR
		{"0101000080000000000000f03f000000000000f03f000000000000f03f", 0, XYZ, POINT, []float64{1, 1, 1}}, // EWKB HDR
		{"00000003e93ff00000000000003ff00000000000003ff0000000000000", 0, XYZ, POINT, []float64{1, 1, 1}}, // WKB XDR
		{"01e9030000000000000000f03f000000000000f03f000000000000f03f", 0, XYZ, POINT, []float64{1, 1, 1}}, // WKB HDR

		// XYZS
		{"00a000000100006c343ff00000000000003ff00000000000003ff0000000000000", 27700, XYZS, POINT, []float64{1, 1, 1}}, // EWKB XDR
		{"01010000a0346c0000000000000000f03f000000000000f03f000000000000f03f", 27700, XYZS, POINT, []float64{1, 1, 1}}, // EWKB HDR
		{"00000003e93ff00000000000003ff00000000000003ff0000000000000", 0, XYZ, POINT, []float64{1, 1, 1}},              // WKB XDR
		{"01e9030000000000000000f03f000000000000f03f000000000000f03f", 0, XYZ, POINT, []float64{1, 1, 1}},              // WKB HDR

		// XYM
		{"00400000013ff00000000000003ff00000000000003ff0000000000000", 0, XYM, POINT, []float64{1, 1, 1}}, // EWKB XDR
		{"0101000040000000000000f03f000000000000f03f000000000000f03f", 0, XYM, POINT, []float64{1, 1, 1}}, // EWKB HDR
		{"00000007d13ff00000000000003ff00000000000003ff0000000000000", 0, XYM, POINT, []float64{1, 1, 1}}, // WKB XDR
		{"01d1070000000000000000f03f000000000000f03f000000000000f03f", 0, XYM, POINT, []float64{1, 1, 1}}, // WKB HDR

		// XYMS
		{"006000000100006c343ff00000000000003ff00000000000003ff0000000000000", 27700, XYMS, POINT, []float64{1, 1, 1}}, // EWKB XDR
		{"0101000060346c0000000000000000f03f000000000000f03f000000000000f03f", 27700, XYMS, POINT, []float64{1, 1, 1}}, // EWKB HDR
		{"00000007d13ff00000000000003ff00000000000003ff0000000000000", 0, XYM, POINT, []float64{1, 1, 1}},              // WKB XDR
		{"01d1070000000000000000f03f000000000000f03f000000000000f03f", 0, XYM, POINT, []float64{1, 1, 1}},              // WKB HDR

		// XYZM
		{"00c00000013ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 0, XYZM, POINT, []float64{1, 1, 1, 1}}, // EWKB XDR
		{"01010000c0000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 0, XYZM, POINT, []float64{1, 1, 1, 1}}, // EWKB HDR
		{"0000000bb93ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 0, XYZM, POINT, []float64{1, 1, 1, 1}}, // WKB XDR
		{"01b90b0000000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 0, XYZM, POINT, []float64{1, 1, 1, 1}}, // WKB HDR

		// XYZM
		{"00e000000100006c343ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 27700, XYZMS, POINT, []float64{1, 1, 1, 1}}, // EWKB XDR
		{"01010000e0346c0000000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 27700, XYZMS, POINT, []float64{1, 1, 1, 1}}, // EWKB HDR
		{"0000000bb93ff00000000000003ff00000000000003ff00000000000003ff0000000000000", 0, XYZM, POINT, []float64{1, 1, 1, 1}},              // WKB XDR
		{"01b90b0000000000000000f03f000000000000f03f000000000000f03f000000000000f03f", 0, XYZM, POINT, []float64{1, 1, 1, 1}},              // WKB HDR
	}

	for _, dataset := range datasets {
		t.Log("Decoding HEX String: DIM = ", dataset.dim, ", Data = ", dataset.data)
		data, err := hex.DecodeString(dataset.data)

		if err != nil {
			t.Fatal("Failed to decode HEX string: err = ", err)
		}

		r := bytes.NewReader(data)

		geom, err := Decode(r)

		if err != nil {
			t.Fatal("Failed to parse geom: err = ", err)
		}

		point := geom.(*Point)

		t.Log(point)

		assert.Equal(t, dataset.dim, point.Dimension(), "Expected dim %v, but got %v ", dataset.dim, geom.Dimension)
		assert.Equal(t, dataset.gtype, point.Type(), "Expected type %v, but got %v ", dataset.gtype, geom.Type)
		assert.Equal(t, dataset.srid, point.Srid(), "Expected srid %v, but got %v ", dataset.srid, geom.Srid)
		assert.Equal(t, len(dataset.points), len(point.Coordinate))

		for idx := 0; idx < len(dataset.points); idx++ {
			assert.InDelta(t, dataset.points[idx], point.Coordinate[idx], 1e-9)
		}
	}
}
