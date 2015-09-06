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

	buf.WriteString(fmt.Sprintf(`package colfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)
`))

	if err := d.appendStructSpec(buf); err != nil {
		return err
	}

	buf.WriteString(fmt.Sprintf(`
func (o *%s) Unmarshal(data []byte) error {
	for i := 0; i < len(data); {
		key := data[i]
		i++
		flag := key&0x80 != 0
		field := key & 0x7f

		switch field {
		default:
			return fmt.Errorf("field %%d unknown", field)
`, d.Name))

	for _, field := range d.Fields {
		s, err := field.switchCase()
		if err != nil {
			return err
		}
		buf.WriteString(s)
	}

	buf.WriteString(`		}
	}
	return nil
}
`)

	f, err := os.Create(fmt.Sprintf("./colfer-%s.go", d.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	buf.WriteTo(f)
	return nil
}

func (o *Object) appendStructSpec(buf *bytes.Buffer) error {
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

func (f *Field) switchCase() (code string, err error) {
	switch f.Type {
	default:
		return "", fmt.Errorf("colfer: type %s unsupported", f.Type)

	case "bool":
		code = fmt.Sprintf(`		case %d:

			o.%s = flag
`, f.No, f.Name)

	case "int32":
		code = fmt.Sprintf(`		case %d:
			x, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			if n < 0 || x&0xFFFFFFFF80000000 != 0 {
				return errors.New("colfer: field %s overflow")
			}

			v := int32(x)
			if flag {
				v = -v
			}
			o.%s = v
`, f.No, f.Name, f.Name)

	case "string":
		code = fmt.Sprintf(`		case %d:
			if flag {
				return errors.New("field %s flag reserved")
			}

			length, n := binary.Uvarint(data[i:])
			if n <= 0 || length&0xFFFFFFFF80000000 != 0 {
				return errors.New("colfer: field %s corrupt size")
			}
			i += n

			to := i + int(length)
			if to > len(data) || to < 0 {
				return io.EOF
			}
			o.%s = string(data[i:to])
			i = to
`, f.No, f.Name, f.Name, f.Name)

	}
	return
}
