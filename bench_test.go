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
//go:generate protoc --gogofaster_out=. -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/gogo/protobuf/protobuf testdata/bench/scheme.proto

var testSet = make([]*bench.Colfer, 1000)

func init()  {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	typ := reflect.TypeOf(testSet).Elem()
	for i := range testSet {
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

		testSet[i] = o
	}
}

var holdData []byte

func BenchmarkEncode(b *testing.B) {
	buf := make([]byte, 512)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		holdData = testSet[i%len(testSet)].Marshal(buf)
	}
}

func BenchmarkEncodeProtoBuf(b *testing.B) {
	buf := make([]byte, 512)
	protoBufSet := make([]*bench.ProtoBuf, len(testSet))
	for i, o := range testSet {
		protoBufSet[i] = &bench.ProtoBuf{
			Key: o.Key,
			Host: o.Host,
			Addr: o.Addr,
			Port: o.Port,
			Size_: o.Size,
			Hash: o.Hash,
			Ratio: o.Ratio,
			Route: o.Route,
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		i, err := protoBufSet[i%len(protoBufSet)].MarshalTo(buf)
		if err != nil {
			b.Error(err)
		}
		holdData = buf[:i]
	}
}

var holdColfer = new(bench.Colfer)

func BenchmarkDecode(b *testing.B) {
	serials := make([][]byte, len(testSet))
	for i, o := range testSet {
		serials[i] = o.Marshal(make([]byte, 1024))
	}

	zero := new(bench.Colfer)
	buf := new(bench.Colfer)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		*buf = *zero	// init
		buf.Unmarshal(serials[i%len(serials)])
		*holdColfer = *buf
	}
}

var holdProtoBuf = new(bench.ProtoBuf)

func BenchmarkDecodeProtoBuf(b *testing.B) {
	serials := make([][]byte, len(testSet))
	for i, o := range testSet {
		p := &bench.ProtoBuf{
			Key: o.Key,
			Host: o.Host,
			Addr: o.Addr,
			Port: o.Port,
			Size_: o.Size,
			Hash: o.Hash,
			Ratio: o.Ratio,
			Route: o.Route,
		}
		var err error
		serials[i], err = p.Marshal()
		if err != nil {
			b.Error(err)
		}
	}

	zero := new(bench.ProtoBuf)
	buf := new(bench.ProtoBuf)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		*buf = *zero	// init
		holdProtoBuf := new(bench.ProtoBuf)
		if err := buf.Unmarshal(serials[i%len(serials)]); err != nil {
			b.Error(err)
		}
		*holdProtoBuf = *buf
	}
}
