package poc

import (
	"encoding/binary"
	"math"
	"math/bits"
	"unsafe"
)

// ColferMax limits serial sizes to 16 MiB.
const ColferMax = 16 * 1024 * 1024

var masks = [...]uint64{
	0,
	0xff,
	0xffff,
	0xffffff,
	0xffffffff,
	0xffffffffff,
	0xffffffffffff,
	0xffffffffffffff,
	0xffffffffffffffff,
}

type Record struct {
	Key   int64
	Host  string
	Port  uint16
	Size  int64
	Hash  uint64
	Ratio float64
	Route bool
}

// MarshalTo encodes o as Colfer into buf. The return is zero when ColferMax was
// reached. Otherwise, the return contains the byte size of the serial written,
// as in serial := buf[:o.MarshalTo(buf)].
func (o *Record) MarshalTo(buf *[ColferMax]byte) int {
	// words of fixed section
	var word0 uint64 = 22 - 1
	var word1 uint64
	var word2 uint64
	var word3 uint64

	// write cursor at variable section
	i := uint64(25)

	// pack Key int64
	v := uint64(o.Key>>63) ^ uint64(o.Key<<1)
	if v < 128 {
		v = v<<1 | 1
	} else {
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56)
		bitCount := uint64(bits.Len64(v))
		extraN := (((bitCount - 1) >> 3) + bitCount) >> 3
		i += extraN
		v >>= uint(extraN)<<3 - 1
		v = (v | 1) << extraN
	}
	word0 |= v << 24

	// pack Host text size
	v = uint64(len(o.Host))
	if v < 128 {
		v = v<<1 | 1
	} else {
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56)
		bitCount := uint64(bits.Len64(v))
		extraN := (((bitCount - 1) >> 3) + bitCount) >> 3
		i += extraN
		v >>= uint(extraN)<<3 - 1
		v = (v | 1) << extraN
	}
	word0 |= v << 32

	// pack Port uint16
	word0 |= uint64(o.Port) << 40

	// pack Size int64
	v = uint64(o.Size>>63) ^ uint64(o.Size<<1)
	if v < 128 {
		v = v<<1 | 1
	} else {
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		p[0], p[1], p[2], p[3], p[4], p[5], p[6], p[7] = byte(v), byte(v>>8), byte(v>>16), byte(v>>24), byte(v>>32), byte(v>>40), byte(v>>48), byte(v>>56)
		bitCount := uint64(bits.Len64(v))
		extraN := (((bitCount - 1) >> 3) + bitCount) >> 3
		i += extraN
		v >>= uint(extraN)<<3 - 1
		v = (v | 1) << extraN
	}
	word0 |= v << 56

	// pack Hash opaque64
	word1 = o.Hash

	binary.LittleEndian.PutUint64(buf[8:], word1)

	// pack Ratio float64
	word2 = math.Float64bits(o.Ratio)

	binary.LittleEndian.PutUint64(buf[16:], word2)

	// pack Route bool
	if o.Route {
		word3 |= 1 << 0
	}

	// write header tail
	buf[24] = byte(word3)

	// copy payloads
	for i <= uint64(len(buf)) {
		p := buf[i:]
		if len(p) < len(o.Host) {
			break
		}
		copy(p, o.Host)
		i += uint64(len(o.Host))

		// finish header
		word0 |= uint64(i)<<17 | 1<<16
		binary.LittleEndian.PutUint64(buf[:], word0)
		return int(i)
	}

	return 0
}

// Unmarshal decodes buf as Colfer into o. It returns either the number of bytes
// read, or zero when an error occurred.
func (o *Record) Unmarshal(buf *[ColferMax]byte) int {
	// words of fixed section
	word0 := binary.LittleEndian.Uint64(buf[:])
	word1 := binary.LittleEndian.Uint64(buf[8:])
	word2 := binary.LittleEndian.Uint64(buf[16:])
	word3 := binary.LittleEndian.Uint64(buf[24:])

	// read cursor at variable section
	i := word0&0xffff + 4

	// unpack variable size
	v := word0 >> 17 & 0x7f
	if word0&(1<<16) == 0 {
		tz := uint64(bits.TrailingZeros64(v|0x80)&7) + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		v |= masks[tz] & (uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 | uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56)
		i += tz
	}
	size := v

	// unpack Key int64
	v = word0 >> 25 & 0x7f
	if word0&(1<<24) == 0 {
		tz := uint64(bits.TrailingZeros64(v|0x80)&7) + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		v |= masks[tz] & (uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 | uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56)
		i += tz
	}
	o.Key = int64(v>>1) ^ -int64(v&1)

	// unpack Host text size
	v = word0 >> 33 & 0x7f
	if word0&(1<<32) == 0 {
		tz := uint64(bits.TrailingZeros64(v|0x80)&7) + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		v |= masks[tz] & (uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 | uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56)
		i += tz
	}
	size_host := v

	// unpack Port uint16
	o.Port = uint16(word0 >> 40)

	// unpack Size int64
	v = word0 >> 57
	if word0&(1<<56) == 0 {
		tz := uint64(bits.TrailingZeros64(v|0x80)&7) + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		p := (*[8]byte)(unsafe.Add(unsafe.Pointer(buf), i))
		v |= masks[tz] & (uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 | uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56)
		i += tz
	}
	o.Size = int64(v>>1) ^ -int64(v&1)

	// unpack Hash opaque64
	o.Hash = word1

	// unpack Ratio float64
	o.Ratio = math.Float64frombits(word2)

	// unpack Route bool
	o.Route = word3&1 != 0

	if l := word0 & 0xffff; l < 22-1 {
		// clear/undo absent fields
		switch l {
		default:
			return 0
		case 1 - 1:
			size_host = 0
			fallthrough
		case 2 - 1:
			o.Port = 0
			fallthrough
		case 4 - 1:
			o.Size = 0
			fallthrough
		case 5 - 1:
			o.Hash = 0
			fallthrough
		case 13 - 1:
			o.Ratio = 0
			fallthrough
		case 21 - 1:
			o.Route = false
		}
	}

	// define serial end
	if size > uint64(len(buf)) {
		return 0
	}
	serial := buf[:size]

	// copy payloads
	offset := size - size_host
	if offset > uint64(len(serial)) {
		return 0
	}
	o.Host = string(serial[offset:])

	return int(size)
}
