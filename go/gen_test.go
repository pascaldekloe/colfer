package testdata

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/pascaldekloe/goe/verify"

	"github.com/pascaldekloe/colfer/go/gen"
)

type golden struct {
	serial string
	object gen.O
}

func newGoldenCases() []*golden {
	return []*golden{
		{"7f", gen.O{}},
		{"007f", gen.O{B: true}},
		{"01017f", gen.O{U32: 1}},
		{"01ff017f", gen.O{U32: math.MaxUint8}},
		{"01ffff037f", gen.O{U32: math.MaxUint16}},
		{"81ffffffff7f", gen.O{U32: math.MaxUint32}},
		{"02017f", gen.O{U64: 1}},
		{"02ff017f", gen.O{U64: math.MaxUint8}},
		{"02ffff037f", gen.O{U64: math.MaxUint16}},
		{"02ffffffff0f7f", gen.O{U64: math.MaxUint32}},
		{"82ffffffffffffffff7f", gen.O{U64: math.MaxUint64}},
		{"03017f", gen.O{I32: 1}},
		{"83017f", gen.O{I32: -1}},
		{"037f7f", gen.O{I32: math.MaxInt8}},
		{"8380017f", gen.O{I32: math.MinInt8}},
		{"03ffff017f", gen.O{I32: math.MaxInt16}},
		{"838080027f", gen.O{I32: math.MinInt16}},
		{"03ffffffff077f", gen.O{I32: math.MaxInt32}},
		{"8380808080087f", gen.O{I32: math.MinInt32}},
		{"04017f", gen.O{I64: 1}},
		{"84017f", gen.O{I64: -1}},
		{"047f7f", gen.O{I64: math.MaxInt8}},
		{"8480017f", gen.O{I64: math.MinInt8}},
		{"04ffff017f", gen.O{I64: math.MaxInt16}},
		{"848080027f", gen.O{I64: math.MinInt16}},
		{"04ffffffff077f", gen.O{I64: math.MaxInt32}},
		{"8480808080087f", gen.O{I64: math.MinInt32}},
		{"04ffffffffffffffff7f7f", gen.O{I64: math.MaxInt64}},
		{"848080808080808080807f", gen.O{I64: math.MinInt64}},
		{"05000000017f", gen.O{F32: math.SmallestNonzeroFloat32}},
		{"057f7fffff7f", gen.O{F32: math.MaxFloat32}},
		{"057fc000007f", gen.O{F32: float32(math.NaN())}},
		{"0600000000000000017f", gen.O{F64: math.SmallestNonzeroFloat64}},
		{"067fefffffffffffff7f", gen.O{F64: math.MaxFloat64}},
		{"067ff80000000000017f", gen.O{F64: math.NaN()}},
		{"0755ef312a2e5da4e77f", gen.O{T: time.Unix(1441739050, 777888999).In(time.UTC)}},
		{"87000007dba8218000000003e87f", gen.O{T: time.Unix(864E10, 1000).In(time.UTC)}},
		{"87fffff82457de8000000003e97f", gen.O{T: time.Unix(-864E10, 1001).In(time.UTC)}},
		{"87ffffffffffffffff2e5da4e77f", gen.O{T: time.Unix(-1, 777888999).In(time.UTC)}},
		{"0801417f", gen.O{S: "A"}},
		{"080261007f", gen.O{S: "a\x00"}},
		{"0809c280e0a080f09080807f", gen.O{S: "\u0080\u0800\U00010000"}},
		{"08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f", gen.O{S: strings.Repeat(" ", 128)}},
		{"0901ff7f", gen.O{A: []byte{math.MaxUint8}}},
		{"090202007f", gen.O{A: []byte{2, 0}}},
		{"09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f", gen.O{A: bytes.Repeat([]byte{9}, 192)}},
		{"0a7f7f", gen.O{O: &gen.O{}}},
		{"0a007f7f", gen.O{O: &gen.O{B: true}}},
		{"0b01007f7f", gen.O{Os: []*gen.O{{B: true}}}},
		{"0b027f7f7f", gen.O{Os: []*gen.O{{}, {}}}},
		{"0c0300016101627f", gen.O{Ss: []string{"", "a", "b"}}},
		{"0d0201000201027f", gen.O{As: [][]byte{[]byte{0}, []byte{1, 2}}}},
		{"0e017f", gen.O{U8: 1}},
		{"0eff7f", gen.O{U8: math.MaxUint8}},
		{"8f017f", gen.O{U16: 1}},
		{"0fffff7f", gen.O{U16: math.MaxUint16}},
		{"1002000000003f8000007f", gen.O{F32s: []float32{0, 1}}},
		{"11014058c000000000007f", gen.O{F64s: []float64{99}}},
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

		got := gen.O{}
		if err := got.UnmarshalBinary(data); err != nil {
			t.Errorf("0x%s: %s", gold.serial, err)
			continue
		}
		verify.Values(t, fmt.Sprintf("0x%s", gold.serial), got, gold.object)
	}
}

func TestUnmarshalEOF(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := hex.DecodeString(gold.serial)
		if err != nil {
			t.Fatal(err)
		}

		for i := range data {
			incomplete := data[:i]
			if _, err := new(gen.O).Unmarshal(incomplete); err != io.EOF {
				t.Errorf("0x%s: got error %T: %q", hex.EncodeToString(incomplete), err, err)
			}
		}
	}
}

func TestUnmarshalSizeMax(t *testing.T) {
	orig := gen.ColferSizeMax
	defer func() {
		gen.ColferSizeMax = orig
	}()

	for _, gold := range newGoldenCases() {
		data, err := hex.DecodeString(gold.serial)
		if err != nil {
			t.Fatal(err)
		}

		for gen.ColferSizeMax = range data {
			// not supported
			if gen.ColferSizeMax == 0 {
				continue
			}

			for i := range data {
				// cutoff on or after max
				if i+1 < gen.ColferSizeMax {
					continue
				}
				part := data[:i+1]

				switch _, err := new(gen.O).Unmarshal(part); err.(type) {
				case gen.ColferMax:
					continue // pass
				case nil:
					t.Errorf("0x%s: no error with ColferSizeMax=%d", hex.EncodeToString(part), gen.ColferSizeMax)
				default:
					t.Errorf("0x%s: got error %T with ColferSizeMax=%d: %q", hex.EncodeToString(part), err, gen.ColferSizeMax, err)
				}
				break
			}
		}
	}
}

// TestFuzzSeed updates the initial input corpus for fuzz testing.
func TestFuzzSeed(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := gold.object.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		if ioutil.WriteFile("../testdata/corpus/seed"+gold.serial, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
