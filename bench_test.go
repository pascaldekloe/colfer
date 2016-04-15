package colfer

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	"github.com/pascaldekloe/colfer/testdata/bench"
)

//go:generate go run ./cmd/colf/main.go -p testdata go testdata/bench/scheme.colf
//go:generate go run ./cmd/colf/main.go -p testdata java testdata/bench/scheme.colf
//go:generate go run ./cmd/colf/main.go -b testdata/bench ecmascript testdata/bench/scheme.colf
//go:generate protoc --gogofaster_out=. -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/gogo/protobuf/protobuf testdata/bench/scheme.proto
//go:generate flatc -o testdata -g testdata/bench/scheme.fbs

var testSet = make([]*bench.Colfer, 1000)

func init() {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	typ := reflect.TypeOf(testSet).Elem()
	for i := range testSet {
		v, ok := quick.Value(typ, rnd)
		if !ok {
			panic("can't generate testdata.Colfer values")
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

// holdData prevents compiler optimization.
var holdData []byte

func BenchmarkEncode(b *testing.B) {
	buf := make([]byte, 512)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		n := testSet[i%len(testSet)].MarshalTo(buf)
		holdData = buf[:n]
	}
}

func BenchmarkEncodeProtoBuf(b *testing.B) {
	buf := make([]byte, 512)
	protoBufSet := make([]*bench.ProtoBuf, len(testSet))
	for i, o := range testSet {
		protoBufSet[i] = &bench.ProtoBuf{
			Key:   o.Key,
			Host:  o.Host,
			Addr:  o.Addr,
			Port:  o.Port,
			Size_: o.Size,
			Hash:  o.Hash,
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

func BenchmarkEncodeFlatBuffers(b *testing.B) {
	builder := flatbuffers.NewBuilder(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		o := testSet[i%len(testSet)]

		builder.Reset()
		host := builder.CreateString(o.Host)
		addr := builder.CreateByteVector(o.Addr)
		bench.FlatBuffersStart(builder)
		bench.FlatBuffersAddKey(builder, o.Key)
		bench.FlatBuffersAddHost(builder, host)
		bench.FlatBuffersAddAddr(builder, addr)
		bench.FlatBuffersAddPort(builder, o.Port)
		bench.FlatBuffersAddSize(builder, o.Size)
		bench.FlatBuffersAddHash(builder, o.Hash)
		bench.FlatBuffersAddRatio(builder, o.Ratio)
		if o.Route {
			bench.FlatBuffersAddRoute(builder, 1)
		} else {
			bench.FlatBuffersAddRoute(builder, 0)
		}
		builder.Finish(bench.FlatBuffersEnd(builder))

		holdData = builder.Bytes[builder.Head():]
	}
}

// holdColfer prevents compiler optimization.
var holdColfer *bench.Colfer

func BenchmarkDecode(b *testing.B) {
	serials := make([][]byte, len(testSet))
	for i, o := range testSet {
		serials[i], _ = o.MarshalBinary()
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		data := serials[i%len(serials)]
		holdColfer = new(bench.Colfer)
		if err := holdColfer.UnmarshalBinary(data); err != nil {
			b.Fatal(err)
		}
	}
}

// holdProtoBuf prevents compiler optimization.
var holdProtoBuf *bench.ProtoBuf

func BenchmarkDecodeProtoBuf(b *testing.B) {
	serials := make([][]byte, len(testSet))
	for i, o := range testSet {
		p := &bench.ProtoBuf{
			Key:   o.Key,
			Host:  o.Host,
			Addr:  o.Addr,
			Port:  o.Port,
			Size_: o.Size,
			Hash:  o.Hash,
			Ratio: o.Ratio,
			Route: o.Route,
		}
		var err error
		serials[i], err = p.Marshal()
		if err != nil {
			b.Error(err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		data := serials[i%len(serials)]
		holdProtoBuf := new(bench.ProtoBuf)
		if err := holdProtoBuf.Unmarshal(data); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDecodeFlatBuf(b *testing.B) {
	serials := make([][]byte, len(testSet))
	for i, o := range testSet {
		builder := flatbuffers.NewBuilder(0)
		host := builder.CreateString(o.Host)
		addr := builder.CreateByteVector(o.Addr)
		bench.FlatBuffersStart(builder)
		bench.FlatBuffersAddKey(builder, o.Key)
		bench.FlatBuffersAddHost(builder, host)
		bench.FlatBuffersAddAddr(builder, addr)
		bench.FlatBuffersAddPort(builder, o.Port)
		bench.FlatBuffersAddSize(builder, o.Size)
		bench.FlatBuffersAddHash(builder, o.Hash)
		bench.FlatBuffersAddRatio(builder, o.Ratio)
		if o.Route {
			bench.FlatBuffersAddRoute(builder, 1)
		} else {
			bench.FlatBuffersAddRoute(builder, 0)
		}
		builder.Finish(bench.FlatBuffersEnd(builder))
		serials[i] = builder.FinishedBytes()
	}

	buf := new(bench.FlatBuffers)

	b.ReportAllocs()
	b.ResetTimer()
	for i := b.N; i != 0; i-- {
		data := serials[i%len(serials)]
		buf.Init(data, flatbuffers.GetUOffsetT(data))

		holdColfer = new(bench.Colfer)
		holdColfer.Key = buf.Key()
		holdColfer.Host = string(buf.Host())
		n := buf.AddrLength()
		holdColfer.Addr = make([]byte, n)
		for i := 0; i < n; i++ {
			holdColfer.Addr[i] = byte(buf.Addr(i))
		}
		holdColfer.Port = buf.Port()
		holdColfer.Size = buf.Size()
		holdColfer.Hash = buf.Hash()
		holdColfer.Ratio = buf.Ratio()
		holdColfer.Route = buf.Route() == 1
	}
}
