package colfer

import (
	"encoding/binary"
	"fmt"
	"io"
)

type tstobj struct {
	b   bool
	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64
	u   uint
	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64
	f32 float32
	f64 float64
	s   string
}

func (o *tstobj) Unmarshal(data []byte) error {
	for i := 0; i < len(data); {
		key := data[i]
		i++
		flag := key&0x80 != 0
		field := key & 0x7f

		switch field {
		default:
			return fmt.Errorf("field %d unknown", field)
		case 0:
			if flag {
				fmt.Errorf("field %d flag reserved", field)
			}
			o.b = true
		case 1:
			x, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			// BUG(ps) Detect int byte size
			if n < 0 || x&0xFFFFFFFF80000000 != 0 {
				return fmt.Errorf("field %d overflow", field)
			}

			v := int(x)
			if flag {
				v = -v
			}
			o.i = v

		case 13:
			if flag {
				fmt.Errorf("field %d flag reserved", field)
			}

			length, n := binary.Uvarint(data[i:])
			if n <= 0 || length&0xFFFFFFFF80000000 != 0 {
				return fmt.Errorf("field %d corrupt size", field)
			}
			i += n

			to := i + int(length)
			if to > len(data) || to < 0 {
				return io.EOF
			}
			o.s = string(data[i:to])
			i = to

		}
	}
	return nil
}
