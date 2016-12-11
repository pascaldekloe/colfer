package bench

import (
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
)

var testData = []*Colfer{
	{Key: 1234567890, Host: "db003lz12", Port: 389, Size: 452, Hash: 0x488b5c2428488918, Ratio: 0.99, Route: true},
	{Key: 1234567891, Host: "localhost", Port: 22, Size: 4096, Hash: 0x243048899c24c824, Ratio: 0.20, Route: false},
	{Key: 1234567892, Host: "kdc.local", Port: 88, Size: 1984, Hash: 0x000048891c24485c, Ratio: 0.06, Route: false},
	{Key: 1234567893, Host: "vhost8.dmz.example.com", Port: 27017, Size: 59741, Hash: 0x5c2408488b9c2489, Ratio: 0.0, Route: true},
}

var protoTestData = []*ProtoBuf{
	{Key: testData[0].Key, Host: testData[0].Host, Port: uint32(testData[0].Port), Size_: testData[0].Size, Hash: testData[0].Hash, Ratio: testData[0].Ratio, Route: testData[0].Route},
	{Key: testData[1].Key, Host: testData[1].Host, Port: uint32(testData[1].Port), Size_: testData[1].Size, Hash: testData[1].Hash, Ratio: testData[1].Ratio, Route: testData[1].Route},
	{Key: testData[2].Key, Host: testData[2].Host, Port: uint32(testData[2].Port), Size_: testData[2].Size, Hash: testData[2].Hash, Ratio: testData[2].Ratio, Route: testData[2].Route},
	{Key: testData[3].Key, Host: testData[3].Host, Port: uint32(testData[3].Port), Size_: testData[3].Size, Hash: testData[3].Hash, Ratio: testData[3].Ratio, Route: testData[3].Route},
}

var colferSerials = make([][]byte, len(testData))
var protoSerials = make([][]byte, len(protoTestData))
var flatSerials = make([][]byte, len(testData))

func init() {
	for i, o := range testData {
		var err error
		colferSerials[i], err = o.MarshalBinary()
		if err != nil {
			panic(err)
		}
	}

	for i, o := range protoTestData {
		var err error
		protoSerials[i], err = o.Marshal()
		if err != nil {
			panic(err)
		}
	}

	for i, o := range testData {
		builder := flatbuffers.NewBuilder(0)
		host := builder.CreateString(o.Host)
		FlatBuffersStart(builder)
		FlatBuffersAddKey(builder, o.Key)
		FlatBuffersAddHost(builder, host)
		FlatBuffersAddPort(builder, o.Port)
		FlatBuffersAddSize(builder, o.Size)
		FlatBuffersAddHash(builder, o.Hash)
		FlatBuffersAddRatio(builder, o.Ratio)
		if o.Route {
			FlatBuffersAddRoute(builder, 1)
		} else {
			FlatBuffersAddRoute(builder, 0)
		}
		builder.Finish(FlatBuffersEnd(builder))
		flatSerials[i] = builder.FinishedBytes()
	}
}

// prevent compiler optimization
var (
	holdSerial       []byte
	holdData         *Colfer
	holdProtoBufData *ProtoBuf
)

func BenchmarkMarshal(b *testing.B) {
	b.Run("colfer", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			var err error
			holdSerial, err = testData[i%len(testData)].MarshalBinary()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("protobuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			var err error
			holdSerial, err = protoTestData[i%len(testData)].Marshal()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("flatbuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := testData[i%len(testData)]

			builder := flatbuffers.NewBuilder(0)
			host := builder.CreateString(o.Host)
			FlatBuffersStart(builder)
			FlatBuffersAddKey(builder, o.Key)
			FlatBuffersAddHost(builder, host)
			FlatBuffersAddPort(builder, o.Port)
			FlatBuffersAddSize(builder, o.Size)
			FlatBuffersAddHash(builder, o.Hash)
			FlatBuffersAddRatio(builder, o.Ratio)
			if o.Route {
				FlatBuffersAddRoute(builder, 1)
			} else {
				FlatBuffersAddRoute(builder, 0)
			}
			builder.Finish(FlatBuffersEnd(builder))

			holdSerial = builder.Bytes[builder.Head():]
		}
	})
}

func BenchmarkUnmarshal(b *testing.B) {
	b.Run("colfer", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := new(Colfer)
			holdData = o

			_, err := o.Unmarshal(colferSerials[i%len(colferSerials)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("protobuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := new(ProtoBuf)
			holdProtoBufData = o

			err := o.Unmarshal(protoSerials[i%len(protoSerials)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("flatbuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := new(Colfer)
			holdData = o

			bytes := flatSerials[i%len(flatSerials)]
			buf := new(FlatBuffers)
			buf.Init(bytes, flatbuffers.GetUOffsetT(bytes))
			o.Key = buf.Key()
			o.Host = string(buf.Host())
			o.Port = buf.Port()
			o.Size = buf.Size()
			o.Hash = buf.Hash()
			o.Ratio = buf.Ratio()
			o.Route = buf.Route() == 1
		}
	})
}

func BenchmarkMarshalReuse(b *testing.B) {
	buf := make([]byte, ColferSizeMax)

	b.Run("colfer", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := testData[i%len(testData)]

			l, err := o.MarshalLen()
			if err != nil {
				b.Fatal(err)
			}

			o.MarshalTo(buf)
			holdSerial = buf[:l]
		}
	})

	b.Run("protobuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := protoTestData[i%len(protoTestData)]

			l := o.Size()
			if l > len(buf) {
				buf = make([]byte, l+100)
			}

			_, err := o.MarshalTo(buf)
			if err != nil {
				b.Fatal(err)
			}
			holdSerial = buf[:l]
		}
	})

	builder := flatbuffers.NewBuilder(0)

	b.Run("flatbuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			o := testData[i%len(testData)]

			builder.Reset()
			host := builder.CreateString(o.Host)
			FlatBuffersStart(builder)
			FlatBuffersAddKey(builder, o.Key)
			FlatBuffersAddHost(builder, host)
			FlatBuffersAddPort(builder, o.Port)
			FlatBuffersAddSize(builder, o.Size)
			FlatBuffersAddHash(builder, o.Hash)
			FlatBuffersAddRatio(builder, o.Ratio)
			if o.Route {
				FlatBuffersAddRoute(builder, 1)
			} else {
				FlatBuffersAddRoute(builder, 0)
			}
			builder.Finish(FlatBuffersEnd(builder))
			holdSerial = builder.Bytes[builder.Head():]
		}
	})
}

func BenchmarkUnmarshalReuse(b *testing.B) {
	holdData = new(Colfer)
	holdProtoBufData = new(ProtoBuf)

	b.Run("colfer", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			*holdData = Colfer{}
			_, err := holdData.Unmarshal(colferSerials[i%len(colferSerials)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("protobuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			*holdProtoBufData = ProtoBuf{}
			err := holdProtoBufData.Unmarshal(protoSerials[i%len(protoSerials)])
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	buf := new(FlatBuffers)

	b.Run("flatbuf", func(b *testing.B) {
		b.ReportAllocs()
		for i := b.N; i > 0; i-- {
			bytes := flatSerials[i%len(flatSerials)]
			buf.Init(bytes, flatbuffers.GetUOffsetT(bytes))
			*holdData = Colfer{}
			holdData.Key = buf.Key()
			holdData.Host = string(buf.Host())
			holdData.Port = buf.Port()
			holdData.Size = buf.Size()
			holdData.Hash = buf.Hash()
			holdData.Ratio = buf.Ratio()
			holdData.Route = buf.Route() == 1
		}
	})
}
