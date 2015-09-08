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
	"fmt"
	"io"
	"math"
)

var (
	ErrMagicMismatch = errors.New("colfer: magic header mismatch")
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

		case "int32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		u := uint32(v)
		if v < 0 {
			u = ^u + 1
			data[i] = 0x%02x
		} else {
			data[i] = 0x%02x
		}
		i++
		i += binary.PutUvarint(data[i:], uint64(u))
	}
`, f.Name, f.Num|0x80, f.Num)

		case "float32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0.0 {
		data[i] = 0x%02x
		i++

		u := math.Float32bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(u>>24), byte(u>>16), byte(u>>8), byte(u)
		i += 4
	}
`, f.Name, f.Num)

		case "text", "binary":
			fmt.Fprintf(buf, `
	if v := o.%s; len(v) != 0 {
		data[i] = 0x%02x
		i++

		i += binary.PutUvarint(data[i:], uint64(len(v)))

		to := i + len(v)
		if to > len(data) {
			panic("TODO(ps) grow for blob")
		}
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
	if data[0] != 0x%02x {
		return ErrMagicMismatch
	}

	for i := 1; i < len(data); {
		key := data[i]
		i++
		field := key & 0x7f

		switch field {
		default:
			return fmt.Errorf("field %%d unknown", field)
`, o.Name, Magic)

	for _, f := range o.Fields {
		if err := f.writeUnmarshalSwitchCase(buf); err != nil {
			return err
		}
	}

	buf.WriteString(`		}
	}
	return nil
}
`)
	return nil
}

func (f *Field) writeUnmarshalSwitchCase(buf *bytes.Buffer) error {
	switch f.Type {
	default:
		return fmt.Errorf("colfer: type %s unsupported", f.Type)

	case "bool":
		fmt.Fprintf(buf, `		case %d:
			o.%s = true
`, f.Num, f.Name)

	case "int32":
		fmt.Fprintf(buf, `		case %d:
			x, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			if key&0x80 != 0 {
				x = ^x + 1
			}
			o.%s = int32(x)
`, f.Num, f.Name)

	case "float32":
		fmt.Fprintf(buf, `		case %d:
			to := i + 4
			if to < 0 || to > len(data) {
				return io.EOF
			}

			x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
			o.%s = math.Float32frombits(x)
			i = to
`, f.Num, f.Name)

	case "text":
		fmt.Fprintf(buf, `		case %d:
			length, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			to := i + int(length)
			if to < 0 || to > len(data) {
				return io.EOF
			}
			o.%s = string(data[i:to])
			i = to
`, f.Num, f.Name)

	case "binary":
		fmt.Fprintf(buf, `		case %d:
			length, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			to := i + int(length)
			if to < 0 || to > len(data) {
				return io.EOF
			}
			v := make([]byte, to-i)
			copy(v, data[i:to])
			o.%s = v
			i = to
`, f.Num, f.Name)

	}
	return nil
}
