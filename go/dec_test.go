package colfer

import (
	"encoding/hex"
	"reflect"
	"testing"
)

//go:generate go run ../cmd/colf/main.go

type tstobj struct {
	b   bool
	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64
	u   uint
	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64
	f32 float32
	f64 float64
	s   string
}

var golden = []struct {
	serial  string
	mapping tstobj
}{
	{"00", tstobj{}},
	{"80", tstobj{b: true}},
	{"0101", tstobj{i: 1}},
	{"8101", tstobj{i: -1}},
	{"0d0141", tstobj{s: "A"}},
}

func TestGoldenDecodes(t *testing.T) {
	for _, gold := range golden {
		data, err := hex.DecodeString(gold.serial)
		if err != nil {
			t.Fatal(err)
		}

		got := new(tstobj)
		if err = got.Unmarshal(data); err != nil {
			t.Errorf("%s: %s", gold.serial, err)
		}
		if !reflect.DeepEqual(got, &gold.mapping) {
			t.Errorf("%s: got %+v, want %+v", gold.serial, *got, gold.mapping)
		}
	}
}
