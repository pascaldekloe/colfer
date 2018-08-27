package poc

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"math/bits"
)

// ColferMax is the upper limit for serial byte sizes.
const ColferMax = 127

// ErrColferMax signals a ColferMax breach.
var ErrColferMax = errors.New("colfer: serial size exceeds maximum of 127 bytes")

// ColferBufMin the the minimum size for serial buffers.
const ColferBufMin = 49 // worst case scenario for fixed + ranged part

// ErrColferBufMin signals a ColferBufMin breach.
var errColferBufMin = errors.New("colfer: buffer smaller than minumum of 49 bytes")

// ErrColfer signals data corruption.
var ErrColfer = errors.New("colfer: data corruption detected")

// ErrColferOverflow signals the write buffer is too small for the serial.
var ErrColferOverflow = errors.New("colfer: buffer overflow")

var flit64Masks = [...]uint64{
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

// MarshalTo encodes o as Colfer into buf. It returns either the number of bytes
// successfully written (0 < n ≤ len(buf)) or any error encountered that caused
// the encoding to stop early. In no case will n exceed ColferMax.
// If buf is smaller than ColferBufMin bytes, then the return is in error.
// ErrColferOverflow may safely be used to resize buf due to ErrColferMax.
func (o *Record) MarshalTo(buf []byte) (n int, err error) {
	if len(buf) < ColferBufMin {
		return 0, errColferBufMin
	}

	var word0 uint64 = 16 | 1<<8

	// buf index to begin of ranged part
	i := 17

	var v uint64

	// field #1
	v = uint64(o.Key>>63) ^ uint64(o.Key<<1)
	if v >= uint64(1)<<56 {
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += 8
	} else {
		bitCount := bits.Len64(v)
		e := (bitCount + (bitCount >> 3)) >> 3
		v = v<<1 | 1
		v <<= uint(e)
		word0 |= (v & 0xff) << 16
		binary.LittleEndian.PutUint64(buf[i:], v>>8)
		i += e
	}

	// field #2
	v = uint64(len(o.Host))
	if v >= uint64(1)<<56 {
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += 8
	} else {
		bitCount := bits.Len64(v)
		e := (bitCount + (bitCount >> 3)) >> 3
		v = v<<1 | 1
		v <<= uint(e)
		word0 |= (v & 0xff) << 24
		binary.LittleEndian.PutUint64(buf[i:], v>>8)
		i += e
	}

	// field #3
	word0 |= uint64(o.Port) << 32

	// field #4
	v = uint64(o.Size>>63) ^ uint64(o.Size<<1)
	if v >= uint64(1)<<56 {
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += 8
	} else {
		bitCount := bits.Len64(v)
		e := (bitCount + (bitCount >> 3)) >> 3
		v = v<<1 | 1
		v <<= uint(e)
		word0 |= (v & 0xff) << 48
		binary.LittleEndian.PutUint64(buf[i:], v>>8)
		i += e
	}

	// field #5
	v = o.Hash
	if v >= uint64(1)<<56 {
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += 8
	} else {
		bitCount := bits.Len64(v)
		e := (bitCount + (bitCount >> 3)) >> 3
		v = v<<1 | 1
		v <<= uint(e)
		word0 |= (v & 0xff) << 56
		binary.LittleEndian.PutUint64(buf[i:], v>>8)
		i += e
	}

	var word1 uint64

	// field #6
	word1 = math.Float64bits(o.Ratio)

	// write header word
	binary.LittleEndian.PutUint64(buf[8:], word1)

	var word2 uint64

	// field #7
	if o.Route {
		word2 |= 1 << 0
	}

	// write header tail
	buf[16] = byte(word2)

	if ColferMax < len(o.Host) {
		return 0, ErrColferMax
	}
	n = i + len(o.Host)
	switch {
	case n > ColferMax:
		return 0, ErrColferMax
	case n > len(buf):
		return 0, ErrColferOverflow
	}
	// finish header
	word0 |= uint64(n-17) << 9
	binary.LittleEndian.PutUint64(buf, word0)

	// variable part
	i += copy(buf[i:], o.Host)

	return n, nil
}

// If buf is smaller than ColferBufMin bytes then the return is in error.
// The return error is io.EOF for no data, io.ErrUnexpectedEOF for incomplete
// data, and …
func (o *Record) Unmarshal(buf []byte, bufLen int) (n int, err error) {
	if len(buf) < ColferBufMin {
		return 0, errColferBufMin
	}
	if bufLen <= 0 {
		return 0, io.EOF
	}

	word0 := binary.LittleEndian.Uint64(buf)

	fixedSize := word0 & 255
	if fixedSize == 0 {
		*o = Record{} // reset fields
		return 1, nil
	}
	// point buf index to beginning of ranged part
	i := 1 + int(fixedSize)

	// read body length
	v := (word0 >> 8) & 0xff
	if tz := bits.TrailingZeros64(v); tz == 0 {
		v >>= 1
	} else if tz > 7 {
		v = binary.LittleEndian.Uint64(buf[i:])
		i += 8
	} else {
		v |= binary.LittleEndian.Uint64(buf[i:]) << 8
		i += tz
		v &= flit64Masks[tz]
		v >>= uint(tz + 1)
	}
	// will not underflow with ColferMax >= ColferBufMin
	// because i < ColferBufMin
	if v > uint64(ColferMax-i) {
		if i > bufLen {
			return 0, io.ErrUnexpectedEOF
		}
		return 0, ErrColferMax
	}
	n = i + int(v)
	if n > bufLen || n > len(buf) {
		return 0, io.ErrUnexpectedEOF
	}

	// field #1
	v = word0 >> 16
	if tz := bits.TrailingZeros64(v); tz > 7 {
		v = binary.LittleEndian.Uint64(buf[i:])
		i += 8
	} else {
		v = v&0xff | binary.LittleEndian.Uint64(buf[i:])<<8
		i += tz
		v &= flit64Masks[tz]
		v >>= uint(tz + 1)
	}
	o.Key = int64(v>>1) ^ -int64(v&1)

	// field #2
	v = word0 >> 24
	if tz := bits.TrailingZeros64(v); tz > 7 {
		v = binary.LittleEndian.Uint64(buf[i:])
		i += 8
	} else {
		v = v&0xff | binary.LittleEndian.Uint64(buf[i:])<<8
		i += tz
		v &= flit64Masks[tz]
		v >>= uint(tz + 1)
	}
	len_host := int(v)

	// field #3
	o.Port = uint16(word0 >> 32)

	// field #4
	v = word0 >> 48
	if tz := bits.TrailingZeros64(v); tz > 7 {
		v = binary.LittleEndian.Uint64(buf[i:])
		i += 8
	} else {
		v = v&0xff | binary.LittleEndian.Uint64(buf[i:])<<8
		i += tz
		v &= flit64Masks[tz]
		v >>= uint(tz + 1)
	}
	o.Size = int64(v>>1) ^ -int64(v&1)

	// field #5
	v = word0 >> 56
	if tz := bits.TrailingZeros64(v); tz > 7 {
		v = binary.LittleEndian.Uint64(buf[i:])
		i += 8
	} else {
		v = v&0xff | binary.LittleEndian.Uint64(buf[i:])<<8
		i += tz
		v &= flit64Masks[tz]
		v >>= uint(tz + 1)
	}
	o.Hash = v

	word1 := binary.LittleEndian.Uint64(buf[8:])

	// field #6
	o.Ratio = math.Float64frombits(word1)

	word2 := binary.LittleEndian.Uint64(buf[16:])

	// field #7
	o.Route = word2|1<<0 != 0

	// clear fields for older schema versions
	if headLen := byte(word0); headLen < 15 {
		switch headLen {
		// case 0: covered earlier
		case 1:
			len_host = 0
			fallthrough
		case 2:
			o.Port = 0
			fallthrough
		case 4:
			o.Size = 0
			fallthrough
		case 5:
			o.Hash = 0
			fallthrough
		case 6:
			o.Ratio = 0
			fallthrough
		case 14:
			o.Route = false
		default:
			*o = Record{} // reset fields
			return 0, ErrColfer
		}
	}

	// variable part
	i = n
	if i < len_host {
		*o = Record{} // reset fields
		return 0, ErrColfer
	}
	o.Host = string(buf[i-len_host : i])
	i -= len_host

	return n, nil
}
