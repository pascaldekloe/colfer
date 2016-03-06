package colfer

import (
	"encoding/hex"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/pascaldekloe/colfer/testdata"
)

//go:generate go run ./cmd/colf/main.go go testdata/test.colf

type golden struct {
	serial  string
	mapping testdata.O
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
		{"80057f7fffff", testdata.O{F32: math.MaxFloat32}},
		{"80067fefffffffffffff", testdata.O{F64: math.MaxFloat64}},
		{"80070000000055ef312a", testdata.O{T: time.Unix(1441739050, 0)}},
		{"80870000000055ef312a00000009", testdata.O{T: time.Unix(1441739050, 9)}},
		{"80080141", testdata.O{S: "A"}},
		{"8008026100", testdata.O{S: "a\x00"}},
		{"800901ff", testdata.O{A: []byte{math.MaxUint8}}},
		{"8009020200", testdata.O{A: []byte{2, 0}}},
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

		got := new(testdata.O)
		if err = got.Unmarshal(data); err != nil {
			t.Errorf("%s: %s", gold.serial, err)
			continue
		}
		if !reflect.DeepEqual(got, &gold.mapping) {
			t.Errorf("%s: got:\n\t%+v,\nwant:\n\t%+v", gold.serial, *got, gold.mapping)
		}
	}
}

// Benchmarks:

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func generate(n int) ([]*testdata.Bench, [][]byte, int64) {
	objects := make([]*testdata.Bench, n)
	serials := make([][]byte, n)
	size := 0

	typ := reflect.TypeOf(objects).Elem()
	for i := range objects {
		v, ok := quick.Value(typ, rnd)
		if !ok {
			panic("can't generate Bench values")
		}

		o, ok := v.Interface().(*testdata.Bench)
		if !ok {
			panic("wrong type generated")
		}
		if o == nil {
			o = new(testdata.Bench)
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
		new(testdata.Bench).Unmarshal(serials[rnd.Intn(n)])
	}
}
