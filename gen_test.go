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
		{"807f", testdata.O{}},
		{"80007f", testdata.O{B: true}},
		{"8001017f", testdata.O{U32: 1}},
		{"8001ffffffff0f7f", testdata.O{U32: math.MaxUint32}},
		{"8002017f", testdata.O{U64: 1}},
		{"8002ffffffffffffffffff017f", testdata.O{U64: math.MaxUint64}},
		{"8003017f", testdata.O{I32: 1}},
		{"8083017f", testdata.O{I32: -1}},
		{"8003ffffffff077f", testdata.O{I32: math.MaxInt32}},
		{"808380808080087f", testdata.O{I32: math.MinInt32}},
		{"8004017f", testdata.O{I64: 1}},
		{"8084017f", testdata.O{I64: -1}},
		{"8004ffffffffffffffff7f7f", testdata.O{I64: math.MaxInt64}},
		{"8084808080808080808080017f", testdata.O{I64: math.MinInt64}},
		{"8005000000017f", testdata.O{F32: math.SmallestNonzeroFloat32}},
		{"80057f7fffff7f", testdata.O{F32: math.MaxFloat32}},
		{"80057fc000007f", testdata.O{F32: float32(math.NaN())}},
		{"800600000000000000017f", testdata.O{F64: math.SmallestNonzeroFloat64}},
		{"80067fefffffffffffff7f", testdata.O{F64: math.MaxFloat64}},
		{"80067ff80000000000017f", testdata.O{F64: math.NaN()}},
		{"80070000000055ef312a7f", testdata.O{T: time.Unix(1441739050, 0)}},
		{"80870000000055ef312a2e5da4e77f", testdata.O{T: time.Unix(1441739050, 777888999)}},
		{"800801417f", testdata.O{S: "A"}},
		{"80080261007f", testdata.O{S: "a\x00"}},
		{"800809c280e0a080f09080807f", testdata.O{S: "\u0080\u0800\U00010000"}},
		{"800901ff7f", testdata.O{A: []byte{math.MaxUint8}}},
		{"80090202007f", testdata.O{A: []byte{2, 0}}},
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
