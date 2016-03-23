package testdata

import (
	"errors"
	"io"
	"math"
	"time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = math.E
var _ = time.RFC3339

var (
	ErrColferStruct   = errors.New("colfer: struct header mismatch")
	ErrColferField    = errors.New("colfer: unknown field header")
	ErrColferOverflow = errors.New("colfer: varint overflow")
)

type O struct {
	B	bool
	U32	uint32
	U64	uint64
	I32	int32
	I64	int64
	F32	float32
	F64	float64
	T	time.Time
	S	string
	A	[]byte
}

// MarshalTo encodes o as Colfer into buf and returns the number of bytes written.
// If the buffer is too small, MrashalTo will panic.
func (o *O) MarshalTo(buf []byte) int {
	if o == nil {
		return 0
	}

	buf[0] = 0x80
	i := 1

	if o.B {
		buf[i] = 0x00
		i++
	}

	if x := o.U32; x != 0 {
		buf[i] = 0x01
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if x := o.U64; x != 0 {
		buf[i] = 0x02
		i++
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.I32; v != 0 {
		buf[i] = 0x03
		i++
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
			buf[i-1] |= 0x80
		}
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.I64; v != 0 {
		buf[i] = 0x04
		i++
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
			buf[i-1] |= 0x80
		}
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
	}

	if v := o.F32; v != 0.0 {
		buf[i] = 0x05
		i++
		x := math.Float32bits(v)
		buf[i], buf[i+1], buf[i+2], buf[i+3] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 4
	}

	if v := o.F64; v != 0.0 {
		buf[i] = 0x06
		i++
		x := math.Float64bits(v)
		buf[i], buf[i+1], buf[i+2], buf[i+3] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		buf[i+4], buf[i+5], buf[i+6], buf[i+7] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 8
	}

	if v := o.T; !v.IsZero() {
		buf[i] = 0x07
		i++
		s, ns := v.Unix(), v.Nanosecond()
		buf[i], buf[i+1], buf[i+2], buf[i+3] = byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32)
		buf[i+4], buf[i+5], buf[i+6], buf[i+7] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
		i += 8
		if ns != 0 {
			buf[i-9] |= 0x80
			buf[i], buf[i+1], buf[i+2], buf[i+3] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 4
		}
	}

	if v := o.S; len(v) != 0 {
		buf[i] = 0x08
		i++
		x := uint(len(v))
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		to := i + len(v)
		copy(buf[i:], v)
		i = to
	}

	if v := o.A; len(v) != 0 {
		buf[i] = 0x09
		i++
		x := uint(len(v))
		for x >= 0x80 {
			buf[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		buf[i] = byte(x)
		i++
		to := i + len(v)
		copy(buf[i:], v)
		i = to
	}

	return i
}

// MarshalSize returns the number of bytes that will hold the Colfer serial for sure.
func (o *O) MarshalSize() int {
	if o == nil {
		return 0
	}

	// BUG(pascaldekloe): MarshalBinary panics on documents larger than 2kB due to the
	// fact that MarshalSize is not implemented yet.
	return 2048
}

// MarshalBinary encodes o as Colfer conform encoding.BinaryMarshaler.
// The error return is always nil.
func (o *O) MarshalBinary() (data []byte, err error) {
	data = make([]byte, o.MarshalSize())
	n := o.MarshalTo(data)
	return data[:n], nil
}

// UnmarshalBinary decodes data as Colfer conform encoding.BinaryUnmarshaler.
// The error return options are io.EOF, ErrColferStruct, ErrColferField and ErrColferOverflow.
func (o *O) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}

	if data[0] != 0x80 {
		return ErrColferStruct
	}
	if len(data) == 1 {
		return nil
	}

	header := data[1]
	field := header & 0x7f
	i := 2

	if field == 0 {
		o.B = true

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 1 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if shift >= 32 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		o.U32 = x

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 2 {
		var x uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint64(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		o.U64 = x

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 3 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if shift >= 32 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		if header&0x80 != 0 {
			x = ^x + 1
		}
		o.I32 = int32(x)

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 4 {
		var x uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint64(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		if header&0x80 != 0 {
			x = ^x + 1
		}
		o.I64 = int64(x)

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 5 {
		to := i + 4
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		o.F32 = math.Float32frombits(x)
		i = to

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 6 {
		to := i + 8
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.F64 = math.Float64frombits(x)
		i = to

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 7 {
		to := i + 8
		var nsec int64
		if header&0x80 == 0 {
			if to < 0 || to > len(data) {
				return io.EOF
			}
		} else {
			to += 4
			if to < 0 || to > len(data) {
				return io.EOF
			}
			nsec = int64(uint(data[i+8])<<24 | uint(data[i+9])<<16 | uint(data[i+10])<<8 | uint(data[i+11]))
		}
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		i = to

		o.T = time.Unix(int64(sec), nsec)

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 8 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if shift >= 32 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		to := i + int(x)
		if to < 0 || to > len(data) {
			return io.EOF
		}
		o.S = string(data[i:to])
		i = to

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	if field == 9 {
		var x uint32
		for shift := uint(0); ; shift += 7 {
			if shift >= 32 {
				return ErrColferOverflow
			}
			b := data[i]
			i++
			x |= (uint32(b) & 0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		length := int(x)
		to := i + length
		if to < 0 || to > len(data) {
			return io.EOF
		}
		v := make([]byte, length)
		copy(v, data[i:to])
		o.A = v
		i = to

		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}

	return ErrColferField
}
