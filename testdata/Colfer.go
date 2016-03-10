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
	ErrStructMismatch = errors.New("colfer: struct header mismatch")
	ErrCorrupt        = errors.New("colfer: data corrupt")
	ErrOverflow       = errors.New("colfer: integer overflow")
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

func (o *O) Marshal(data []byte) []byte {
	data[0] = 0x80
	i := 1

	if o.B {
		data[i] = 0x00
		i++
	}

	if x := o.U32; x != 0 {
		data[i] = 0x01
		i++
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}

	if x := o.U64; x != 0 {
		data[i] = 0x02
		i++
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}

	if v := o.I32; v != 0 {
		data[i] = 0x03
		i++
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
			data[i-1] |= 0x80
		}
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}

	if v := o.I64; v != 0 {
		data[i] = 0x04
		i++
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
			data[i-1] |= 0x80
		}
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}

	if v := o.F32; v != 0.0 {
		data[i] = 0x05
		i++
		x := math.Float32bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 4
	}

	if v := o.F64; v != 0.0 {
		data[i] = 0x06
		i++
		x := math.Float64bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		data[i+4], data[i+5], data[i+6], data[i+7] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 8
	}

	if v := o.T; !v.IsZero() {
		data[i] = 0x07
		i++
		s, ns := v.Unix(), v.Nanosecond()
		data[i], data[i+1], data[i+2], data[i+3] = byte(s>>56), byte(s>>48), byte(s>>40), byte(s>>32)
		data[i+4], data[i+5], data[i+6], data[i+7] = byte(s>>24), byte(s>>16), byte(s>>8), byte(s)
		i += 8
		if ns != 0 {
			data[i-9] |= 0x80
			data[i], data[i+1], data[i+2], data[i+3] = byte(ns>>24), byte(ns>>16), byte(ns>>8), byte(ns)
			i += 4
		}
	}

	if v := o.S; len(v) != 0 {
		data[i] = 0x08
		i++
		x := uint(len(v))
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
		to := i + len(v)
		copy(data[i:], v)
		i = to
	}

	if v := o.A; len(v) != 0 {
		data[i] = 0x09
		i++
		x := uint(len(v))
		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
		to := i + len(v)
		copy(data[i:], v)
		i = to
	}

	return data[:i]
}

func (o *O) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	if data[0] != 0x80 {
		return ErrStructMismatch
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
				return ErrOverflow
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
				return ErrOverflow
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
				return ErrOverflow
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
				return ErrOverflow
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
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		i += 8

		var nsec int64
		if header&0x80 != 0 {
			v := uint(data[i])<<24 | uint(data[i+1])<<16 | uint(data[i+2])<<8 | uint(data[i+3])
			i += 4
			nsec = int64(v)
		}

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
				return ErrOverflow
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
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
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
				return ErrOverflow
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
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
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

	return ErrCorrupt
}
