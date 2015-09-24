package colfer

import (
	"os"
	"text/template"
)

// Generate writes the code into file "Colfer.go".
func Generate(pkg *Package) error {
	t := template.New("go-code").Delims("<:", ":>")
	template.Must(t.Parse(goCode))
	template.Must(t.New("marshal").Parse(goMarshal))
	template.Must(t.New("marshal-field").Parse(goMarshalField))
	template.Must(t.New("unmarshal").Parse(goUnmarshal))
	template.Must(t.New("unmarshal-field").Parse(goUnmarshalField))

	f, err := os.Create("Colfer.go")
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, pkg)
}

const goCode = `package <:.Name:>

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
        ErrMagicMismatch = errors.New("colfer: magic header mismatch")
        ErrCorrupt       = errors.New("colfer: data corrupt")
)

<:range .Objects:>
type <:.Name:> struct {
<:range .Fields:>	<:.Name:>	<:if eq .Type "timestamp":>time.Time<:else if eq .Type "text":>string<:else if eq .Type "binary":>[]byte<:else:><:.Type:><:end:>
<:end:>}
<:end:>

<:range .Objects:>
<:template "marshal" .:>
<:template "unmarshal" .:><:end:>`

const goMarshal = `func (o *<:.Name:>) Marshal(data []byte) []byte {
	data[0] = 0x80
	i := 1
<:range .Fields:><:template "marshal-field" .:><:end:>
	return data[:i]
}
`

const goMarshalField = `<:$fieldDecl := printf "data[i] = 0x%02x" .Num :><:if eq .Type "bool":>
	if o.<:.Name:> == true {
		<:$fieldDecl:>
		i++
	}
<:else if eq .Type "uint32":>
	if v := o.<:.Name:>; v != 0 {
		<:$fieldDecl:>
		i++
		for v >= 0x80 {
			data[i] = byte(v) | 0x80
			v >>= 7
			i++
		}
		data[i] = byte(v)
		i++
	}
<:else if eq .Type "uint64":>
	if v := o.<:.Name:>; v != 0 {
		<:$fieldDecl:>
		i++
		for v >= 0x80 {
			data[i] = byte(v) | 0x80
			v >>= 7
			i++
		}
		data[i] = byte(v)
		i++
	}
<:else if eq .Type "int32":>
	if v := o.<:.Name:>; v != 0 {
		x := uint32(v)
		if v < 0 {
			x = ^x + 1
			<:$fieldDecl:> | 0x80
		} else {
			<:$fieldDecl:>
		}
		i++
		for x >= 0x80 {
			data[i] = byte(v) | 0x80
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}
<:else if eq .Type "int64":>
	if v := o.<:.Name:>; v != 0 {
		x := uint64(v)
		if v < 0 {
			x = ^x + 1
			<:$fieldDecl:> | 0x80
		} else {
			<:$fieldDecl:>
		}
		i++
		for x >= 0x80 {
			data[i] = byte(v) | 0x80
			x >>= 7
			i++
		}
		data[i] = byte(x)
		i++
	}
<:else if eq .Type "float32":>
	if v := o.<:.Name:>; v != 0.0 {
		<:$fieldDecl:>
		x := math.Float32bits(v)
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 5
	}
<:else if eq .Type "float64":>
	if v := o.<:.Name:>; v != 0.0 {
		<:$fieldDecl:>
		x := math.Float64bits(v)
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32)
		data[i+5], data[i+6], data[i+7], data[i+8] = byte(x>>24), byte(x>>16), byte(x>>8), byte(x)
		i += 9
	}
<:else if eq .Type "timestamp":>
	if v := o.<:.Name:>; !v.IsZero() {
		<:$fieldDecl:>
		sec, nsec := v.Unix(), v.Nanosecond()
		data[i+1], data[i+2], data[i+3], data[i+4] = byte(sec>>56), byte(sec>>48), byte(sec>>40), byte(sec>>32)
		data[i+5], data[i+6], data[i+7], data[i+8] = byte(sec>>24), byte(sec>>16), byte(sec>>8), byte(sec)
		if nsec != 0 {
			data[i] |= 0x80
			data[i+9], data[i+10], data[i+11], data[i+12] = byte(nsec>>24), byte(nsec>>16), byte(nsec>>8), byte(nsec)
			i += 4
		}
		i += 9
	}
<:else if eq .Type "text" "binary":>
	if v := o.<:.Name:>; len(v) != 0 {
		<:$fieldDecl:>
		i++
		length := uint(len(v))
		for length >= 0x80 {
			data[i] = byte(length) | 0x80
			length >>= 7
			i++
		}
		data[i] = byte(length)
		i++
		to := i + len(v)
		copy(data[i:], v)
		i = to
	}
<:end:>`

const goUnmarshal = `
func (o *<:.Name:>) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return io.EOF
	}
	if data[0] != 0x80 {
		return ErrMagicMismatch
	}

	if len(data) == 1 {
		return nil
	}
	key := data[1]
	field := key & 0x7f
	i := 2
<:range .Fields:><:template "unmarshal-field" .:><:end:>
	return ErrCorrupt
}
`

const goUnmarshalField = `
	if field < <:.Num:> {
		return ErrCorrupt
	}
	if field == <:.Num:> {
<:if eq .Type "bool":>		o.<:.Name:> = true
<:else if eq .Type "uint32":>
		var x uint32
		for shift := uint(0); shift <= 25; shift += 7 {
			b := data[i]
			i++
			x |= uint32(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		o.<:.Name:> = x
<:else if eq .Type "uint64":>
		var x uint64
		for shift := uint(0); shift <= 57; shift += 7 {
			b := data[i]
			i++
			x |= uint64(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		o.<:.Name:> = x
<:else if eq .Type "int32":>
		var x uint32
		for shift := uint(0); shift <= 25; shift += 7 {
			b := data[i]
			i++
			x |= uint32(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		if key&0x80 != 0 {
			x = ^x + 1
		}
		o.<:.Name:> = int32(x)
<:else if eq .Type "int64":>
		var x uint64
		for shift := uint(0); shift <= 57; shift += 7 {
			b := data[i]
			i++
			x |= uint64(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		if key&0x80 != 0 {
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
		var nsec uint64
		if key&0x80 != 0 {
			nsec = uint64(data[i])<<24 | uint64(data[i+1])<<16 | uint64(data[i+2])<<8 | uint64(data[i+3])
			i += 4
		}
		o.<:.Name:> = time.Unix(int64(sec), int64(nsec))
<:else if eq .Type "text":>
		var length uint
		for shift := uint(0); shift <= 57; shift += 7 {
			b := data[i]
			i++
			length |= uint(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		to := i + int(length)
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
			return io.EOF
		}
		o.<:.Name:> = string(data[i:to])
		i = to
<:else if eq .Type "binary":>
		var length uint
		for shift := uint(0); shift <= 57; shift += 7 {
			b := data[i]
			i++
			length |= uint(b&0x7f) << shift
			if b < 0x80 {
				break
			}
			if i == len(data) {
				return io.EOF
			}
		}
		to := i + int(length)
		if to < 0 {
			return ErrCorrupt
		}
		if to > len(data) {
			return io.EOF
		}
		v := make([]byte, to-i)
		copy(v, data[i:to])
		o.<:.Name:> = v
		i = to
<:end:>
		if i == len(data) {
			return nil
		}
		key = data[i]
		field = key & 0x7f
		i++
	}
`
