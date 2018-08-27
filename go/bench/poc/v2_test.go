package poc

import (
	"bytes"
	"reflect"
	"testing"
)

var benchData = []byte{
	16,        // fixed size
	11<<1 | 1, // ranged + variable size FLIT
	2,         // key FLIT fixed
	8<<1 | 1,  // host size FLIT fixed
	254, 253,  // port fixed
	2,                              // size FLIT fixed
	2,                              // hash FLIT fixed
	10, 11, 12, 13, 14, 15, 16, 17, // ratio fixed
	1,                              // route fixed
	8,                              // key FLIT ranged
	16,                             // size FLIT ranged
	16,                             // hash FLIT ranged
	0, 21, 22, 23, 24, 25, 26, 255, // host variable
}

var benchObj = &Record{Key: 256, Host: "\x00\x15\x16\x17\x18\x19\x1a\xff", Port: 0xfdfe, Size: 512, Hash: 0x400, Ratio: 1.694714631965086e-226, Route: true}

func BenchmarkMarshalTo(b *testing.B) {
	buf := make([]byte, ColferMax)

	o := *benchObj
	for i := 0; i < b.N; i++ {
		n, err := o.MarshalTo(buf)
		if err != nil {
			b.Fatal("marshal error:", err)
		}
		if n != len(benchData) {
			b.Fatalf("wrote %d bytes, want %d", n, len(benchData))
		}
	}

	if got := buf[:len(benchData)]; !bytes.Equal(got, benchData) {
		b.Errorf("got %x, want %x", got, benchData)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	buf := make([]byte, ColferMax)
	copy(buf, benchData)

	var o Record
	for i := 0; i < b.N; i++ {
		n, err := o.Unmarshal(buf, len(benchData))
		if err != nil {
			b.Fatal(err)
		}
		if n != len(benchData) {
			b.Fatalf("read %d out of %d", n, len(benchData))
		}
	}

	if !reflect.DeepEqual(&o, benchObj) {
		b.Errorf("got %#v, want %#v", o, benchObj)
	}
}
