package colfer

import (
	"encoding/hex"
	"reflect"
	"testing"
)

//go:generate go run ../cmd/colf/main.go

var golden = []struct {
	serial  string
	mapping tstobj
}{
	{"", tstobj{}},
	{"80", tstobj{b: true}},
	{"0101", tstobj{i32: 1}},
	{"8101", tstobj{i32: -1}},
	{"024008f5c3", tstobj{f32: 2.14}},
	{"030141", tstobj{s: "A"}},
}

func TestGoldenEncodes(t *testing.T) {
	for _, gold := range golden {
		got := hex.EncodeToString(gold.mapping.Marshal(make([]byte, 1000)))
		if got != gold.serial {
			t.Errorf("Got 0x%s, want 0x%s", got, gold.serial)
		}
	}
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
