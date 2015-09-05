package colfer

import (
	"fmt"
	"os"
	"reflect"
)

type Object struct {
	Name   string
	Fields []*Field
}

func (d *Object) Generate() error {
	f, err := os.Create(fmt.Sprintf("./colfer-%s.go", d.Name))
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf(`package colfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)
	
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

	if err != nil {
		return err
	}

	for _, field := range d.Fields {
		s, err := field.switchCase()
		if err != nil {
			return err
		}
		if _, err = f.WriteString(s); err != nil {
			return err
		}
	}

	f.WriteString(`		}
	}
	return nil
}`)

	return nil
}

type Field struct {
	No   int
	Name string
	Kind reflect.Kind
}

func (f *Field) switchCase() (code string, err error) {
	switch f.Kind {
	default:
		return "", fmt.Errorf("colfer: kind %s unsupported", f.Kind)

	case reflect.Bool:
		code = fmt.Sprintf(`		case %d:

			o.%s = flag
`, f.No, f.Name)

	case reflect.Int:
		code = fmt.Sprintf(`		case %d:
			x, n := binary.Uvarint(data[i:])
			if n == 0 {
				return io.EOF
			}
			i += n

			// BUG(ps) Detect int byte size
			if n < 0 || x&0xFFFFFFFF80000000 != 0 {
				return errors.New("colfer: field %s overflow")
			}

			v := int(x)
			if flag {
				v = -v
			}
			o.%s = v
`, f.No, f.Name, f.Name)

	case reflect.String:
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
