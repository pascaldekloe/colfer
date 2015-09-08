package colfer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func (o *Object) Generate() error {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, `package %s

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"time"
)

var (
	ErrMagicMismatch = errors.New("colfer: magic header mismatch")
	ErrCorrupt       = errors.New("colfer: data corrupt")
)
`, o.Package)

	if err := o.writeTypeCode(buf); err != nil {
		return err
	}

	o.writeMarshalCode(buf)
	o.writeUnmarshalCode(buf)

	f, err := os.Create(fmt.Sprintf("./colf_%s.go", o.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	return err
}

func (o *Object) writeTypeCode(buf *bytes.Buffer) error {
	fields := make([]*ast.Field, len(o.Fields))
	for i, f := range o.Fields {
		t := f.Type
		switch t {
		case "text":
			t = "string"
		case "binary":
			t = "[]byte"
		case "timestamp":
			t = "time.Time"
		}

		fields[i] = &ast.Field{
			Names: []*ast.Ident{
				{Name: f.Name},
			},
			Type: &ast.Ident{
				Name: t,
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`colfer:\"%d\"`", f.Num),
			},
		}
	}

	spec := &ast.TypeSpec{
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: fields,
			},
		},
		Name: &ast.Ident{
			Name: o.Name,
			Obj: &ast.Object{
				Kind: 3,
				Name: o.Name,
			},
		},
	}

	buf.WriteString("\ntype ")
	err := printer.Fprint(buf, token.NewFileSet(), spec)
	buf.WriteByte('\n')
	return err
}

func (o *Object) writeMarshalCode(buf *bytes.Buffer) {
	fmt.Fprintf(buf, `
func (o *%s) Marshal(data []byte) []byte {
	data[0] = 0x%02x
	i := 1
`, o.Name, Magic)

	for _, f := range o.Fields {
		switch f.Type {
		case "bool":
			fmt.Fprintf(buf, `
	if o.%s == true {
		data[i] = 0x%02x
		i++
	}
`, f.Name, f.Num)

		case "uint32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		data[i] = 0x%02x
		i++
		i += binary.PutUvarint(data[i:], uint64(v))
	}
`, f.Name, f.Num)

		case "uint64":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		data[i] = 0x%02x
		i++
		i += binary.PutUvarint(data[i:], v)
	}
`, f.Name, f.Num)

		case "int32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
			data[i] = 0x%02x
		} else {
			data[i] = 0x%02x
		}
		i++
		i += binary.PutUvarint(data[i:], uint64(x))
	}
`, f.Name, f.Num|0x80, f.Num)

		case "int64":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
			data[i] = 0x%02x
		} else {
			data[i] = 0x%02x
		}
		i++
		i += binary.PutUvarint(data[i:], x)
	}
`, f.Name, f.Num|0x80, f.Num)

		case "float32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0.0 {
		x := math.Float32bits(v)
		data[i] = 0x%02x
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	}
`, f.Name, f.Num)

		case "float64":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0.0 {
		x := math.Float64bits(v)
		data[i] = 0x%02x
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		data[i+5], data[i+6], data[i+7], data[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	}
`, f.Name, f.Num)

		case "timestamp":
			fmt.Fprintf(buf, `
	if v := o.%s; !v.IsZero() {
		sec, nsec := v.Unix(), v.Nanosecond()
		data[i] = 0x%02x
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(sec>>56), byte(sec>>48), byte(sec>>40), byte(sec>>32)
		data[i+5], data[i+6], data[i+7], data[i+8] = byte(sec>>24), byte(sec>>16), byte(sec>>8), byte(sec)
		if nsec != 0 {
			data[i] |= 0x80
			data[i+9], data[i+10], data[i+11], data[i+12] = byte(nsec>>24), byte(nsec>>16), byte(nsec>>8), byte(nsec)
			i += 4
		}
		i += 9
	}
`, f.Name, f.Num)

		case "text", "binary":
			fmt.Fprintf(buf, `
	if v := o.%s; len(v) != 0 {
		data[i] = 0x%02x
		i++
		i += binary.PutUvarint(data[i:], uint64(len(v)))
		to := i + len(v)
		copy(data[i:], v)
		i = to
	}
`, f.Name, f.Num)

		}
	}

	buf.WriteString(`
	return data[:i]
}
`)
}

func (o *Object) writeUnmarshalCode(buf *bytes.Buffer) error {
	fmt.Fprintf(buf, `
func (o *%s) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	if data[0] != 0x%02x {
		return ErrMagicMismatch
	}

	if len(data) < 2 {
		return nil
	}
	key := data[1]
	field := key & 0x7f
	i := 2

`, o.Name, Magic)

	for _, f := range o.Fields {
		if err := f.writeUnmarshalCode(buf); err != nil {
			return err
		}
	}

	buf.WriteString(`	return ErrCorrupt
}
`)
	return nil
}

func (f *Field) writeUnmarshalCode(buf *bytes.Buffer) error {
	fmt.Fprintf(buf, `	if field < %d {
		return ErrCorrupt
	}
	if field == %d {`, f.Num, f.Num)

	switch f.Type {
	default:
		return fmt.Errorf("colfer: type %s unsupported", f.Type)

	case "bool":
		fmt.Fprintf(buf, `
		o.%s = true
`, f.Name)

	case "int32":
		fmt.Fprintf(buf, `
		x, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		if key&0x80 != 0 {
			x = ^x + 1
		}
		o.%s = int32(x)
`, f.Name)

	case "int64":
		fmt.Fprintf(buf, `
		x, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		if key&0x80 != 0 {
			x = ^x + 1
		}
		o.%s = int64(x)
`, f.Name)

	case "uint32":
		fmt.Fprintf(buf, `
		x, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		o.%s = uint32(x)
`, f.Name)

	case "uint64":
		fmt.Fprintf(buf, `
		x, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		o.%s = x
`, f.Name)

	case "float32":
		fmt.Fprintf(buf, `
		to := i + 4
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		o.%s = math.Float32frombits(x)
		i = to
`, f.Name)

	case "float64":
		fmt.Fprintf(buf, `
		to := i + 8
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.%s = math.Float64frombits(x)
		i = to
`, f.Name)

	case "timestamp":
		fmt.Fprintf(buf, `
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		i += 8
		var nsec uint64
		if key&0x80 != 0 {
			nsec = uint64(data[i])<<24 | uint64(data[i+1])<<16 | uint64(data[i+2])<<8 | uint64(data[i+3])
			i += 4
		}
		o.%s = time.Unix(int64(sec), int64(nsec))
`, f.Name)

	case "text":
		fmt.Fprintf(buf, `
		length, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		to := i + int(length)
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
			return io.EOF
		}
		o.%s = string(data[i:to])
		i = to
`, f.Name)

	case "binary":
		fmt.Fprintf(buf, `
		length, n := binary.Uvarint(data[i:])
		if n <= 0 {
			if n == 0 {
				return io.EOF
			}
			return ErrCorrupt
		}
		i += n
		to := i + int(length)
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
			return io.EOF
		}
		v := make([]byte, to-i)
		copy(v, data[i:to])
		o.%s = v
		i = to
`, f.Name)

	}

	buf.WriteString(`		if i == len(data) {
			return nil
		}
		key = data[i]
		field = key & 0x7f
		i++
	}

`)
	return nil
}
