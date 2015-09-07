package colfer

import (
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

//go:generate go run ../cmd/colf/main.go

var golden = []struct {
	serial  string
	mapping tstobj
}{
	{"", tstobj{}},
	{"80", tstobj{B: true}},
	{"0101", tstobj{I32: 1}},
	{"8101", tstobj{I32: -1}},
	{"024008f5c3", tstobj{F32: 2.14}},
	{"030141", tstobj{S: "A"}},
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

func generate(n int) ([]*tstobj, [][]byte, int64) {
	objects := make([]*tstobj, n)
	serials := make([][]byte, n)
	size := 0

	r := rand.New(rand.NewSource(8))
	typ := reflect.TypeOf(objects).Elem()
	for i := range objects {
		v, ok := quick.Value(typ, r)
		if !ok {
			panic("Can't generate tstobj")
		}

		o, ok := v.Interface().(*tstobj)
		if o == nil {
			o = new(tstobj)
		}

		if !ok {
			panic("Wrong type generated")
		}
		b := o.Marshal(make([]byte, 1000))

		objects[i], serials[i] = o, b
		size += len(b)
	}

	return objects, serials, int64(size / n)
}

func BenchmarkEncode(b *testing.B) {
	n := 1000
	objects, _, avgSize := generate(n)
	buf := make([]byte, 1000)

	b.SetBytes(avgSize)
	b.ReportAllocs()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		objects[rand.Intn(n)].Marshal(buf)
	}
}

func BenchmarkDecode(b *testing.B) {
	n := 1000
	_, serials, avgSize := generate(n)

	b.SetBytes(avgSize)
	b.ReportAllocs()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		o := new(tstobj)
		o.Unmarshal(serials[rand.Intn(n)])
	}
}
