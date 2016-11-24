package colfer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
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
		{"7f", testdata.O{}},
		{"007f", testdata.O{B: true}},
		{"01017f", testdata.O{U32: 1}},
		{"01ff017f", testdata.O{U32: math.MaxUint8}},
		{"01ffff037f", testdata.O{U32: math.MaxUint16}},
		{"81ffffffff7f", testdata.O{U32: math.MaxUint32}},
		{"02017f", testdata.O{U64: 1}},
		{"02ff017f", testdata.O{U64: math.MaxUint8}},
		{"02ffff037f", testdata.O{U64: math.MaxUint16}},
		{"02ffffffff0f7f", testdata.O{U64: math.MaxUint32}},
		{"82ffffffffffffffff7f", testdata.O{U64: math.MaxUint64}},
		{"03017f", testdata.O{I32: 1}},
		{"83017f", testdata.O{I32: -1}},
		{"037f7f", testdata.O{I32: math.MaxInt8}},
		{"8380017f", testdata.O{I32: math.MinInt8}},
		{"03ffff017f", testdata.O{I32: math.MaxInt16}},
		{"838080027f", testdata.O{I32: math.MinInt16}},
		{"03ffffffff077f", testdata.O{I32: math.MaxInt32}},
		{"8380808080087f", testdata.O{I32: math.MinInt32}},
		{"04017f", testdata.O{I64: 1}},
		{"84017f", testdata.O{I64: -1}},
		{"047f7f", testdata.O{I64: math.MaxInt8}},
		{"8480017f", testdata.O{I64: math.MinInt8}},
		{"04ffff017f", testdata.O{I64: math.MaxInt16}},
		{"848080027f", testdata.O{I64: math.MinInt16}},
		{"04ffffffff077f", testdata.O{I64: math.MaxInt32}},
		{"8480808080087f", testdata.O{I64: math.MinInt32}},
		{"04ffffffffffffffff7f7f", testdata.O{I64: math.MaxInt64}},
		{"848080808080808080807f", testdata.O{I64: math.MinInt64}},
		{"05000000017f", testdata.O{F32: math.SmallestNonzeroFloat32}},
		{"057f7fffff7f", testdata.O{F32: math.MaxFloat32}},
		{"057fc000007f", testdata.O{F32: float32(math.NaN())}},
		{"0600000000000000017f", testdata.O{F64: math.SmallestNonzeroFloat64}},
		{"067fefffffffffffff7f", testdata.O{F64: math.MaxFloat64}},
		{"067ff80000000000017f", testdata.O{F64: math.NaN()}},
		{"0755ef312a2e5da4e77f", testdata.O{T: time.Unix(1441739050, 777888999).In(time.UTC)}},
		{"870000000100000000000000007f", testdata.O{T: time.Unix(math.MaxUint32+1, 0).In(time.UTC)}},
		{"87ffffffffffffffff2e5da4e77f", testdata.O{T: time.Unix(-1, 777888999).In(time.UTC)}},
		{"87fffffff14f443f00000000007f", testdata.O{T: time.Unix(-63094636800, 0).In(time.UTC)}},
		{"0801417f", testdata.O{S: "A"}},
		{"080261007f", testdata.O{S: "a\x00"}},
		{"0809c280e0a080f09080807f", testdata.O{S: "\u0080\u0800\U00010000"}},
		{"08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f", testdata.O{S: strings.Repeat(" ", 128)}},
		{"0901ff7f", testdata.O{A: []byte{math.MaxUint8}}},
		{"090202007f", testdata.O{A: []byte{2, 0}}},
		{"09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f", testdata.O{A: bytes.Repeat([]byte{9}, 192)}},
		{"0a7f7f", testdata.O{O: &testdata.O{}}},
		{"0a007f7f", testdata.O{O: &testdata.O{B: true}}},
		{"0b01007f7f", testdata.O{Os: []*testdata.O{{B: true}}}},
		{"0b027f7f7f", testdata.O{Os: []*testdata.O{{}, {}}}},
		{"0c0300016101627f", testdata.O{Ss: []string{"", "a", "b"}}},
		{"0d0201000201027f", testdata.O{As: [][]byte{[]byte{0}, []byte{1, 2}}}},
		{"0e017f", testdata.O{U8: 1}},
		{"0eff7f", testdata.O{U8: math.MaxUint8}},
		{"8f017f", testdata.O{U16: 1}},
		{"0fffff7f", testdata.O{U16: math.MaxUint16}},
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

// TestFuzzSeed updates the initial input corpus for fuzz testing.
func TestFuzzSeed(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := gold.object.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		if ioutil.WriteFile("testdata/corpus/seed"+gold.serial, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
