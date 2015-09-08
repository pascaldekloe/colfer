package colfer

import (
	"encoding/hex"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

//go:generate go run ../cmd/colf/main.go go

type golden struct {
	serial  string
	mapping TstObj
}

func newGoldenCases() []*golden {
	return []*golden{
		{"80", TstObj{}},
		{"8000", TstObj{B: true}},
		{"800101", TstObj{I32: 1}},
		{"808101", TstObj{I32: -1}},
		{"80024008f5c3", TstObj{F32: 2.14}},
		{"80030000000055ef312a", TstObj{T: time.Unix(1441739050, 0)}},
		{"80830000000055ef312a00000009", TstObj{T: time.Unix(1441739050, 9)}},
		{"80040141", TstObj{S: "A"}},
		{"8005020100", TstObj{A: []byte{1, 0}}},
	}
}

func TestGoldenEncodes(t *testing.T) {
	for _, gold := range newGoldenCases() {
		got := hex.EncodeToString(gold.mapping.Marshal(make([]byte, 100)))
		if got != gold.serial {
			t.Errorf("Got 0x%s, want 0x%s", got, gold.serial)
		}
	}
}

func TestGoldenDecodes(t *testing.T) {
	for _, gold := range newGoldenCases() {
		data, err := hex.DecodeString(gold.serial)
		if err != nil {
			t.Fatal(err)
		}

		got := new(TstObj)
		if err = got.Unmarshal(data); err != nil {
			t.Errorf("%s: %s", gold.serial, err)
		}
		if !reflect.DeepEqual(got, &gold.mapping) {
			t.Errorf("%s: got %+v, want %+v", gold.serial, *got, gold.mapping)
		}
	}
}

func generate(n int) ([]*TstObj, [][]byte, int64) {
	objects := make([]*TstObj, n)
	serials := make([][]byte, n)
	size := 0

	typ := reflect.TypeOf(objects).Elem()
	for i := range objects {
		v, ok := quick.Value(typ, rnd)
		if !ok {
			panic("Can't generate TstObj")
		}

		o, ok := v.Interface().(*TstObj)
		if !ok {
			panic("Wrong type generated")
		}
		if o == nil {
			o = new(TstObj)
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
		objects[rnd.Intn(n)].Marshal(buf)
	}
}

func BenchmarkDecode(b *testing.B) {
	n := 1000
	_, serials, avgSize := generate(n)

	b.SetBytes(avgSize)
	b.ReportAllocs()
	b.ResetTimer()

	for i := b.N; i != 0; i-- {
		new(TstObj).Unmarshal(serials[rnd.Intn(n)])
	}
}
