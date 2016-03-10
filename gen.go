package colfer

import (
	"os"
	"path/filepath"
	"text/template"
)

// Generate writes the code into file "Colfer.go".
func Generate(basedir string, structs []*Struct) error {
	pkgT := template.New("go-header").Delims("<:", ":>")
	template.Must(pkgT.Parse(goPackage))

	t := template.New("go-code").Delims("<:", ":>")
	template.Must(t.Parse(goCode))
	template.Must(t.New("marshal-field").Parse(goMarshalField))
	template.Must(t.New("marshal-fieldDecl").Parse(goMarshalFieldDecl))
	template.Must(t.New("marshal-varint").Parse(goMarshalVarint))
	template.Must(t.New("unmarshal-field").Parse(goUnmarshalField))
	template.Must(t.New("unmarshal-varint32").Parse(goUnmarshalVarint32))
	template.Must(t.New("unmarshal-varint64").Parse(goUnmarshalVarint64))

	pkgFiles := make(map[string]*os.File)

	for _, s := range structs {
		f, ok := pkgFiles[s.Pkg.Name]
		if !ok {
			var err error
			f, err = os.Create(filepath.Join(basedir, s.Pkg.Name, "Colfer.go"))
			if err != nil {
				return err
			}
			defer f.Close()

			pkgFiles[s.Pkg.Name] = f
			if err = pkgT.Execute(f, s.Pkg); err != nil {
				return err
			}
		}
		if err := t.Execute(f, s); err != nil {
			return err
		}
	}
	return nil
}

const goPackage = `package <:.Name:>

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

`

const goCode = `type <:.Name:> struct {
<:range .Fields:>	<:.Name:>	<:if eq .Type "timestamp":>time.Time<:else if eq .Type "text":>string<:else if eq .Type "binary":>[]byte<:else:><:.Type:><:end:>
<:end:>}

func (o *<:.Name:>) Marshal(data []byte) []byte {
	data[0] = 0x80
	i := 1
<:range .Fields:><:template "marshal-field" .:><:end:>
	return data[:i]
}

func (o *<:.Name:>) Unmarshal(data []byte) error {
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
<:range .Fields:><:template "unmarshal-field" .:><:end:>
	return ErrCorrupt
}
`

const goMarshalFieldDecl = `		data[i] = <:printf "0x%02x" .Index:>
		i++`

const goMarshalField = `<:if eq .Type "bool":>
	if o.<:.Name:> {
<:template "marshal-fieldDecl" .:>
	}
<:else if eq .Type "uint32":>
	if x := o.<:.Name:>; x != 0 {
<:template "marshal-fieldDecl" .:>
<:template "marshal-varint":>
	}
<:else if eq .Type "uint64":>
	if x := o.<:.Name:>; x != 0 {
<:template "marshal-fieldDecl" .:>
<:template "marshal-varint":>
	}
<:else if eq .Type "int32":>
	if v := o.<:.Name:>; v != 0 {
<:template "marshal-fieldDecl" .:>
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
			data[i-1] |= 0x80
		}
<:template "marshal-varint":>
	}
<:else if eq .Type "int64":>
	if v := o.<:.Name:>; v != 0 {
<:template "marshal-fieldDecl" .:>
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
			data[i-1] |= 0x80
		}
<:template "marshal-varint":>
	}
<:else if eq .Type "float32":>
	if v := o.<:.Name:>; v != 0.0 {
<:template "marshal-fieldDecl" .:>
		x := math.Float32bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 4
	}
<:else if eq .Type "float64":>
	if v := o.<:.Name:>; v != 0.0 {
<:template "marshal-fieldDecl" .:>
		x := math.Float64bits(v)
		data[i], data[i+1], data[i+2], data[i+3] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		data[i+4], data[i+5], data[i+6], data[i+7] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 8
	}
<:else if eq .Type "timestamp":>
	if v := o.<:.Name:>; !v.IsZero() {
<:template "marshal-fieldDecl" .:>
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
<:else if eq .Type "text" "binary":>
	if v := o.<:.Name:>; len(v) != 0 {
<:template "marshal-fieldDecl" .:>
		x := uint(len(v))
<:template "marshal-varint":>
		to := i + len(v)
		copy(data[i:], v)
		i = to
	}
<:end:>`

const goMarshalVarint = `		for x >= 0x80 {
			data[i] = byte(x | 0x80)
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++`

const goUnmarshalField = `
	if field == <:.Index:> {<:if eq .Type "bool":>
		o.<:.Name:> = true
<:else if eq .Type "uint32":>
<:template "unmarshal-varint32":>
		o.<:.Name:> = x
<:else if eq .Type "uint64":>
<:template "unmarshal-varint64":>
		o.<:.Name:> = x
<:else if eq .Type "int32":>
<:template "unmarshal-varint32":>
		if header&0x80 != 0 {
			x = ^x + 1
		}
		o.<:.Name:> = int32(x)
<:else if eq .Type "int64":>
<:template "unmarshal-varint64":>
		if header&0x80 != 0 {
			x = ^x + 1
		}
		o.<:.Name:> = int64(x)
<:else if eq .Type "float32":>
		to := i + 4
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint32(data[i])<<24 | uint32(data[i+1])<<16 | uint32(data[i+2])<<8 | uint32(data[i+3])
		o.<:.Name:> = math.Float32frombits(x)
		i = to
<:else if eq .Type "float64":>
		to := i + 8
		if to < 0 || to > len(data) {
			return io.EOF
		}
		x := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		x |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		o.<:.Name:> = math.Float64frombits(x)
		i = to
<:else if eq .Type "timestamp":>
		sec := uint64(data[i])<<56 | uint64(data[i+1])<<48 | uint64(data[i+2])<<40 | uint64(data[i+3])<<32
		sec |= uint64(data[i+4])<<24 | uint64(data[i+5])<<16 | uint64(data[i+6])<<8 | uint64(data[i+7])
		i += 8

		var nsec int64
		if header&0x80 != 0 {
			v := uint(data[i])<<24 | uint(data[i+1])<<16 | uint(data[i+2])<<8 | uint(data[i+3])
			i += 4
			nsec = int64(v)
		}

		o.<:.Name:> = time.Unix(int64(sec), nsec)
<:else if eq .Type "text":>
<:template "unmarshal-varint32":>
		to := i + int(x)
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
			return io.EOF
		}
		o.<:.Name:> = string(data[i:to])
		i = to
<:else if eq .Type "binary":>
<:template "unmarshal-varint32":>
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
		o.<:.Name:> = v
		i = to
<:end:>
		if i == len(data) {
			return nil
		}
		header = data[i]
		field = header & 0x7f
		i++
	}
`

const goUnmarshalVarint32 = `		var x uint32
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
		}`

const goUnmarshalVarint64 = `		var x uint64
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
		}`
