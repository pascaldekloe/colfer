package colfer

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/pascaldekloe/colfer/testdata"
	"github.com/pascaldekloe/goe/verify"
)

//go:generate go run ./cmd/colf/main.go go testdata/test.colf

type golden struct {
	serial string
	object testdata.O
}

func newGoldenCases() []*golden {
	return []*golden{
		{"80", testdata.O{}},
		{"8000", testdata.O{B: true}},
		{"800101", testdata.O{U32: 1}},
		{"8001ffffffff0f", testdata.O{U32: math.MaxUint32}},
		{"800201", testdata.O{U64: 1}},
		{"8002ffffffffffffffffff01", testdata.O{U64: math.MaxUint64}},
		{"800301", testdata.O{I32: 1}},
		{"808301", testdata.O{I32: -1}},
		{"8003ffffffff07", testdata.O{I32: math.MaxInt32}},
		{"80838080808008", testdata.O{I32: math.MinInt32}},
		{"800401", testdata.O{I64: 1}},
		{"808401", testdata.O{I64: -1}},
		{"8004ffffffffffffffff7f", testdata.O{I64: math.MaxInt64}},
		{"808480808080808080808001", testdata.O{I64: math.MinInt64}},
		{"800500000001", testdata.O{F32: math.SmallestNonzeroFloat32}},
		{"80057f7fffff", testdata.O{F32: math.MaxFloat32}},
		{"80060000000000000001", testdata.O{F64: math.SmallestNonzeroFloat64}},
		{"80067fefffffffffffff", testdata.O{F64: math.MaxFloat64}},
		{"80070000000055ef312a", testdata.O{T: time.Unix(1441739050, 0)}},
		{"80870000000055ef312a2e5da4e7", testdata.O{T: time.Unix(1441739050, 777888999)}},
		{"80080141", testdata.O{S: "A"}},
		{"8008026100", testdata.O{S: "a\x00"}},
		{"800809c280e0a080f0908080", testdata.O{S: "\u0080\u0800\U00010000"}},
		{"800901ff", testdata.O{A: []byte{math.MaxUint8}}},
		{"8009020200", testdata.O{A: []byte{2, 0}}},
	}
}

func TestMarshal(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := gold.object.MarshalBinary()
		if err != nil {
			t.Errorf("0x%s: %s", gold.serial, err)
			continue
		}
		if got := hex.EncodeToString(data); got != gold.serial {
			t.Errorf("Got 0x%s, want 0x%s", got, gold.serial)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := hex.DecodeString(gold.serial)
		if err != nil {
			t.Fatal(err)
		}

		got := testdata.O{}
		if err := got.UnmarshalBinary(data); err != nil {
			t.Errorf("0x%s: %s", gold.serial, err)
			continue
		}
		verify.Values(t, fmt.Sprintf("0x%s", gold.serial), got, gold.object)
	}
}
