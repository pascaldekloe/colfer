package colfer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

func (d *Object) Generate() error {
	buf := new(bytes.Buffer)

	buf.WriteString(`package colfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)
`)

	if err := d.writeTypeCode(buf); err != nil {
		return err
	}

	d.writeMarshalCode(buf)
	d.writeUnmarshalCode(buf)

	f, err := os.Create(fmt.Sprintf("./colfer-%s.go", d.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	buf.WriteTo(f)
	return nil
}

func (o *Object) writeTypeCode(buf *bytes.Buffer) error {
	fields := make([]*ast.Field, len(o.Fields))
	for i, f := range o.Fields {
		fields[i] = &ast.Field{
			Names: []*ast.Ident{
				{Name: f.Name},
			},
			Type: &ast.Ident{
				Name: f.Type,
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`colfer:\"%d\"`", f.No),
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
	buf.WriteString(`
func (o *tstobj) Marshal(data []byte) []byte {
	i := 0
`)

	for _, f := range o.Fields {
		switch f.Type {
		case "bool":
			fmt.Fprintf(buf, `
	if o.%s == true {
		data[i] = 0x80
		i++
	}
`, f.Name)

		case "int32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0 {
		u := uint32(v)
		if v < 0 {
			u = ^u + 1
			data[i] = 0x81
		} else {
			data[i] = 0x01
		}
		i++
		i += binary.PutUvarint(data[i:], uint64(u))
	}
`, f.Name)

		case "float32":
			fmt.Fprintf(buf, `
	if v := o.%s; v != 0.0 {
		data[i] = 0x02
		i++

		u := math.Float32bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)
		i += 4
	}
`, f.Name)

		case "string":
			fmt.Fprintf(buf, `
	if v := o.%s; v != "" {
		data[i] = 0x03
		i++

		i += binary.PutUvarint(data[i:], uint64(len(v)))

		to := i + len(v)
		if to > len(data) {
			panic("TODO(ps) grow for blob")
		}
		copy(data[i:], v)
		i = to
	}
`, f.Name)
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
	for i := 0; i < len(data); {
		key := data[i]
		i++
		flag := key&0x80 != 0
		field := key & 0x7f

		switch field {
		default:
			return fmt.Errorf("field %%d unknown", field)
`, o.Name)

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
			o.%s = flag
`, f.No, f.Name)

	case "int32":
		fmt.Fprintf(buf, `		case %d:
			x, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			if flag {
				x = ^x + 1
			}
			o.%s = int32(x)
`, f.No, f.Name)

	case "float32":
		fmt.Fprintf(buf, `		case %d:
			to := i + 4
			if to < 0 || to > len(data) {
				return io.EOF
			}

			x := uint32(data[i]) << 24 | uint32(data[i+1]) << 16 | uint32(data[i+2]) << 8 | uint32(data[i+3])
			o.%s = math.Float32frombits(x)
			i = to
`, f.No, f.Name)

	case "string":
		fmt.Fprintf(buf, `		case %d:
			if flag {
				return errors.New("colfer: blob field flag reserved")
			}

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
`, f.No, f.Name)

	}
	return nil
}
