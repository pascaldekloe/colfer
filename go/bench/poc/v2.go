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

// ErrColfer signals data corruption or schema mismatch.
var ErrColfer = errors.New("colfer: incompatible data")

// ErrColferOverflow signals the write buffer is too small for the serial.
var ErrColferOverflow = errors.New("colfer: buffer overflow")

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

// MarshalTo encodes o as Colfer into buf. It returns either the number of bytes
// successfully written (0 < n ≤ len(buf)) or any error encountered that caused
// the encoding to stop early. In no case will n exceed ColferMax.
// ErrColferOverflow may safely be used to resize buf due to ErrColferMax.
func (o *Record) MarshalTo(buf *[ColferMax]byte) (n int, err error) {
	if o == nil {
		buf[0], buf[1] = 0, 0
		return 2, nil
	}

	var word0 uint64 = 22

	// buf index to variable part
	i := 25

	var v uint64

	// field #1
	v = uint64(o.Key>>63) ^ uint64(o.Key<<1)
	if v < 128 {
		v = v<<1 | 1
	} else {
		bitCount := bits.Len64(v)
		e := (((bitCount - 1) >> 3) + bitCount) >> 3
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += e
		v >>= uint(e)<<3 - 1
		v = (v | 1) << uint(e)
	}
	word0 |= (v & 0xff) << 24

	// field #2
	v = uint64(len(o.Host))
	if v < 128 {
		v = v<<1 | 1
	} else {
		bitCount := bits.Len64(v)
		e := (((bitCount - 1) >> 3) + bitCount) >> 3
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += e
		v >>= uint(e)<<3 - 1
		v = (v | 1) << uint(e)
	}
	word0 |= (v & 0xff) << 32

	// field #3
	word0 |= uint64(o.Port) << 40

	// field #4
	v = uint64(o.Size>>63) ^ uint64(o.Size<<1)
	if v < 128 {
		v = v<<1 | 1
	} else {
		bitCount := bits.Len64(v)
		e := (((bitCount - 1) >> 3) + bitCount) >> 3
		binary.LittleEndian.PutUint64(buf[i:], v)
		i += e
		v >>= uint(e)<<3 - 1
		v = (v | 1) << uint(e)
	}
	word0 |= (v & 0xff) << 56

	var word1 uint64

	// field #5
	word1 = o.Hash

	binary.LittleEndian.PutUint64(buf[8:], word1)
	var word2 uint64

	// field #6
	word2 = math.Float64bits(o.Ratio)

	binary.LittleEndian.PutUint64(buf[16:], word2)
	var word3 uint64

	// field #7
	if o.Route {
		word3 |= 1 << 0
	}

	// write header tail
	buf[24] = byte(word3)

	// determine serial size
	n = i
	n += len(o.Host)
	if uint(n) > ColferMax {
		return i, ErrColferMax
	}

	// finish header
	word0 |= uint64(n-25)<<17 | 1<<16
	binary.LittleEndian.PutUint64(buf[:], word0)

	// write variable part
	if n > len(buf) {
		return i, ErrColferOverflow
	}
	i += copy(buf[i:], o.Host)

	return n, nil
}

// Unmarshal decodes buf as Colfer. BufLen limits the number of bytes.
// The return error is io.EOF for no data, io.ErrUnexpectedEOF for incomplete
// data, ErrColferMax for size protection and ErrColfer for corrupted data.
func (o *Record) Unmarshal(buf *[ColferMax]byte, bufLen int) (n int, err error) {
	word0 := binary.LittleEndian.Uint64(buf[:8])

	fixedSize := uint16(word0)
	n = int(fixedSize) + 3
	i := n // buf index at variable component

	// variable size
	v := word0 >> 17 & 0x7f
	if word0&(1<<16) == 0 {
		tz := bits.TrailingZeros64(v|0x80)&7 + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		v |= binary.LittleEndian.Uint64(buf[i:]) & masks[tz]
		i += tz
	}
	n += int(v)

	// check boundaries
	switch {
	case v > ColferMax, n > ColferMax:
		if i <= bufLen {
			return i, ErrColferMax
		}
		fallthrough
	case n > bufLen, n > len(buf):
		if bufLen <= 0 {
			return 0, io.EOF
		}
		return 0, io.ErrUnexpectedEOF
	}

	// field #1
	v = word0 >> 25 & 0x7f
	if word0&(1<<24) == 0 {
		tz := bits.TrailingZeros64(v|0x80)&7 + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		v |= binary.LittleEndian.Uint64(buf[i:]) & masks[tz]
		i += tz
	}
	o.Key = int64(v>>1) ^ -int64(v&1)

	// field #2
	v = word0 >> 33 & 0x7f
	if word0&(1<<32) == 0 {
		tz := bits.TrailingZeros64(v|0x80)&7 + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		v |= binary.LittleEndian.Uint64(buf[i:]) & masks[tz]
		i += tz

		if v > ColferMax {
			return i, ErrColfer
		}
	}
	len_host := int(v)

	// field #3
	o.Port = uint16(word0 >> 40)

	// field #4
	v = word0 >> 57 & 0x7f
	if word0&(1<<56) == 0 {
		tz := bits.TrailingZeros64(v|0x80)&7 + 1
		v = v << uint(tz<<3-tz) &^ masks[tz]
		v |= binary.LittleEndian.Uint64(buf[i:]) & masks[tz]
		i += tz
	}
	o.Size = int64(v>>1) ^ -int64(v&1)

	// next word
	word1 := binary.LittleEndian.Uint64(buf[8:])

	// field #5
	o.Hash = word1

	// next word
	word2 := binary.LittleEndian.Uint64(buf[16:])

	// field #6
	o.Ratio = math.Float64frombits(word2)

	// next word
	word3 := binary.LittleEndian.Uint64(buf[24:])

	// field #7
	o.Route = word3&1 != 0

	if fixedSize < 22 {
		// clear/undo fields
		switch fixedSize {
		default:
			return 0, ErrColfer
		case 0:
			o.Key = 0
			fallthrough
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
		case 13:
			o.Ratio = 0
			fallthrough
		case 21:
			o.Route = false
		}
		i = int(fixedSize) + 2
	}

	// variable part
	offset_host := n - len_host
	if offset_host < i {
		return 0, ErrColfer
	}
	o.Host = string(buf[offset_host:n])

	return n, nil
}
