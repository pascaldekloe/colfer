package poc

import (
	"log"
	"reflect"
	"testing"
)

var TestData = [4]*Record{
	{Key: 1234567890, Host: "db003lz12", Port: 389, Size: 452, Hash: 0x488b5c2428488918, Ratio: 0.99, Route: true},
	{Key: 1234567891, Host: "localhost", Port: 22, Size: 4096, Hash: 0x243048899c24c824, Ratio: 0.20, Route: false},
	{Key: 1234567892, Host: "kdc.local", Port: 88, Size: 1984, Hash: 0x000048891c24485c, Ratio: 0.06, Route: false},
	{Key: 1234567893, Host: "vhost8.dmz.example.com", Port: 27017, Size: 59741, Hash: 0x5c2408488b9c2489, Ratio: 0.0, Route: true},
}

var SerialBytes [4]*[ColferMax]byte
var SerialSizes [4]int

func init() {
	for i, o := range TestData {
		SerialBytes[i] = &[ColferMax]byte{}
		SerialSizes[i] = o.MarshalTo(SerialBytes[i])
		if SerialSizes[i] == 0 {
			log.Fatalf("test record %d marshal error", i)
		}
	}
}

func TestRoundtrip(t *testing.T) {
	for i, bytes := range SerialBytes {
		var o Record
		n := o.Unmarshal(bytes)
		if n != SerialSizes[i] {
			t.Errorf("test record %d parsed %d bytes, want %d", i, n, SerialSizes[i])
		}
		if !reflect.DeepEqual(&o, TestData[i]) {
			t.Errorf("test record %d got %#v, want %#v", i, &o, TestData[i])
		}
	}
}

// prevents compiler optimisation
var R Record
var Buf [ColferMax]byte

func BenchmarkMarshalTo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if TestData[i&3].MarshalTo(&Buf) == 0 {
			b.Fatal("marshal error")
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if R.Unmarshal(SerialBytes[i&3]) == 0 {
			b.Fatal("unmarshal error")
		}
	}
}
