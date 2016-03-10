package colfer

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/pascaldekloe/colfer/testdata/bench"
)

//go:generate go run ./cmd/colf/main.go -p testdata go testdata/bench/scheme.colf
//go:generate go run ./cmd/colf/main.go -p testdata java testdata/bench/scheme.colf

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func generate(n int) ([]*bench.Colfer, [][]byte, int64) {
	objects := make([]*bench.Colfer, n)
	serials := make([][]byte, n)
	size := 0

	typ := reflect.TypeOf(objects).Elem()
	for i := range objects {
		v, ok := quick.Value(typ, rnd)
		if !ok {
			panic("can't generate Bench values")
		}

		o, ok := v.Interface().(*bench.Colfer)
		if !ok {
			panic("wrong type generated")
		}
		if o == nil {
			o = new(bench.Colfer)
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
		new(bench.Colfer).Unmarshal(serials[rnd.Intn(n)])
	}
}
