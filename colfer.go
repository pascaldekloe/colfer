// Package colfer provides schema definitions.
package colfer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"unicode"
)

// datatypes holds all supported names.
var datatypes = map[string]struct{}{
	"bool": {},

	"int8":    {},
	"uint8":   {},
	"opaque8": {},

	"int16":    {},
	"uint16":   {},
	"opaque16": {},

	"uint32":   {},
	"int32":    {},
	"opaque32": {},

	"int64":    {},
	"uint64":   {},
	"opaque64": {},

	"float32":   {},
	"float64":   {},
	"timestamp": {},

	"text": {},
}

type Packages []*Package

func (p Packages) Len() int           { return len(p) }
func (p Packages) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p Packages) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// HasFloat returns whether any of the packages has one or more floating-point fields.
func (p Packages) HasFloat() bool {
	for _, o := range p {
		if o.HasFloat() {
			return true
		}
	}
	return false
}

// HasTimestamp returns whether any of the packages has one or more timestamp fields.
func (p Packages) HasTimestamp() bool {
	for _, o := range p {
		if o.HasTimestamp() {
			return true
		}
	}
	return false
}

// Package is a named definition bundle.
type Package struct {
	// Name is the identification token.
	Name string
	// NameNative is the language specific Name.
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Structs are the type definitions.
	Structs []*Struct
	// SchemaFiles are the source filenames.
	SchemaFiles []string
	// SuperClass is the fully qualified path.
	SuperClass string
	// SuperClassNative is the language specific SuperClass.
	SuperClassNative string
	// Interfaces are the fully qualified paths.
	Interfaces []string
	// InterfaceNatives are the language specific Interfaces.
	InterfaceNatives []string
	// CodeSnippet is helpful in book-keeping functionality.
	CodeSnippet string
}

// DocText returns the documentation lines prefixed with ident.
func (p *Package) DocText(ident string) string {
	return docText(p.Docs, ident)
}

// SchemaFileList returns a listing text.
func (p *Package) SchemaFileList() string {
	switch len(p.SchemaFiles) {
	case 0:
		return ""
	case 1:
		return p.SchemaFiles[0]
	default:
		return strings.Join(p.SchemaFiles[1:], ", ") + " and " + p.SchemaFiles[0]
	}
}

// Refs returns all direct references sorted by name.
func (p *Package) Refs() Packages {
	found := make(map[*Package]struct{})
	for _, t := range p.Structs {
		for _, f := range t.Fields {
			if f.TypeRef != nil && f.TypeRef.Pkg != p {
				found[f.TypeRef.Pkg] = struct{}{}
			}
		}
	}

	var refs Packages
	for r := range found {
		refs = append(refs, r)
	}
	sort.Sort(refs)
	return refs
}

// HasFloat returns whether p has one or more floating-point fields.
func (p *Package) HasFloat() bool {
	for _, t := range p.Structs {
		if t.HasFloat() {
			return true
		}
	}
	return false
}

// HasTimestamp returns whether p has one or more timestamp fields.
func (p *Package) HasTimestamp() bool {
	for _, t := range p.Structs {
		if t.HasTimestamp() {
			return true
		}
	}
	return false
}

// HasList returns whether p has one or more list fields.
func (p *Package) HasList() bool {
	for _, t := range p.Structs {
		if t.HasList() {
			return true
		}
	}
	return false
}

// Struct is a data structure definition.
type Struct struct {
	Pkg *Package
	// Name is the identification token.
	Name string
	// NameNative is the language specific Name.
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Fields are the elements in order of appearance.
	Fields []*Field
	// SchemaFile is the source filename.
	SchemaFile string
	// TagAdd has optional source code additions.
	TagAdd []string

	// FixedSize is the encoding space in the fixed section.
	FixedSize int

	// OverflowMax is the upper boundary for encoding space in the overflow
	// section.
	OverflowMax int

	// SizeMin is the smallest encoding size possible (with only the first
	// field present).
	SizeMin int
}

// FieldsReversed returns .Fields in reversed order.
func (t *Struct) FieldsReversed() []*Field {
	fields := append(([]*Field)(nil), t.Fields...)
	slices.Reverse(fields)
	return fields
}

func (t *Struct) SetFixedPositions() {
	t.FixedSize = 3 // header
	t.OverflowMax = 0
	boolCount := 0
	for _, f := range t.Fields {
		f.FixedIndex = t.FixedSize

		// The boolean index is required before
		// the following size calculations!
		if f.Type == "bool" {
			f.BoolIndex = boolCount
			boolCount++
		}

		n := max(1, f.ElementCount)
		t.FixedSize += n * f.TypeFixedSize()
		t.OverflowMax += n * f.TypeOverflowMax()
	}

	t.SizeMin = 3 + t.Fields[0].TypeFixedSize()*max(1, t.Fields[0].ElementCount)
}

// FixedWordIndices returns the index of each 64-bit word filled by the fixed
// section.
func (t *Struct) FixedWordIndices() []int {
	indices := make([]int, t.FixedSize/8)
	for i := range indices {
		indices[i] = i
	}
	return indices
}

// FixedWordRemainderSize returns the octet count.
func (t *Struct) FixedWordRemainderSize() int {
	return t.FixedSize % 8
}

func (t *Struct) FixedWordRemainderByteIndices() []int {
	indices := make([]int, t.FixedWordRemainderSize())
	for i := range indices {
		indices[i] = i
	}
	return indices
}

// DocText returns the documentation lines prefixed with ident.
func (t *Struct) DocText(indent string) string {
	return docText(t.Docs, indent)
}

// String returns the qualified name.
func (t *Struct) String() string {
	return fmt.Sprintf("%s.%s", t.Pkg.Name, t.Name)
}

func (t *Struct) HasPayloadSection() bool {
	return t.HasText() || t.HasList()
}

// HasBool returns whether t has one or more boolean fields.
func (t *Struct) HasBool() bool {
	for _, f := range t.Fields {
		if f.Type == "bool" {
			return true
		}
	}
	return false
}

// HasFloat returns whether t has one or more floating-point fields.
func (t *Struct) HasFloat() bool {
	for _, f := range t.Fields {
		if f.Type == "float32" || f.Type == "float64" {
			return true
		}
	}
	return false
}

// HasText returns whether t has one or more text fields.
func (t *Struct) HasText() bool {
	for _, f := range t.Fields {
		if f.Type == "text" {
			return true
		}
	}
	return false
}

// HasTimestamp returns whether s has one or more timestamp fields.
func (t *Struct) HasTimestamp() bool {
	for _, f := range t.Fields {
		if f.Type == "timestamp" {
			return true
		}
	}
	return false
}

// HasList returns whether s has one or more list fields.
func (t *Struct) HasList() bool {
	for _, f := range t.Fields {
		if f.TypeList {
			return true
		}
	}
	return false
}

// Field is a Struct member definition.
type Field struct {
	// Struct is the parent.
	Struct *Struct
	// Index is the Struct.Fields position.
	Index int
	// Name is the identification token.
	Name string
	// NameNative is the language specific Name.
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Type is the datatype.
	Type string
	// TypeNative is the language specific Type.
	TypeNative string
	// TypeRef is the Colfer data structure reference.
	TypeRef *Struct

	// ElementCount is the fixed-array size.
	ElementCount int
	// TypeList flags whether the datatype is a list.
	TypeList bool

	// TagAdd has optional source code additions.
	TagAdd []string

	// FixedIndex is the position of the first byte in the serial.
	FixedIndex int
	// BoolIndex is the position in the bit field.
	BoolIndex int
}

// WordIndex locates fixed-data in 64-bit words.
func (f *Field) WordIndex() int {
	return f.FixedIndex / 8
}

// NextWordIndex is relevent when the data overlaps two words.
func (f *Field) NextWordIndex() int {
	return f.WordIndex() + 1
}

// WordShift is the bit position of fixed-data at WordIndex.
func (f *Field) WordShift() int {
	return (f.FixedIndex & 7) * 8
}

// TypeFixedSize returns its space in the fixed section.
func (f *Field) TypeFixedSize() (octets int) {
	if f.TypeList {
		return 1
	}

	switch f.Type {
	case "opaque16":
		return 2
	case "opaque32", "float32":
		return 4
	case "opaque64", "float64", "timestamp":
		return 8

	case "bool":
		if f.FirstInBitField() {
			return 1
		}
		return 0 // travels for free
	}

	return 1
}

// TypeOverflowMax returns the upper boundary for its space in the overflow
// section.
func (f *Field) TypeOverflowMax() (octets int) {
	if f.Type == "text" || f.TypeList {
		// This limit would be an error though.
		return 8
	}

	switch f.Type {
	case "uint16", "int16", "uint32", "int32":
		// This limit would be an error though.
		return 8
	case "uint64", "int64":
		return 8
	}

	return 0
}

// Elements returns the field definition for each array item.
func (f *Field) Elements() []Field {
	size := f.TypeFixedSize()
	if f.ElementCount == 0 || size < 1 {
		return nil
	}

	elements := make([]Field, f.ElementCount)
	for i := range elements {
		elements[i] = *f // copy
		elements[i].Index = f.Index + i
		elements[i].FixedIndex = f.FixedIndex + i*size
	}
	return elements
}

// FirstInBitField returns whether this boolean is the first (of up to 8 in
// total) in the byte encoding.
func (f *Field) FirstInBitField() bool {
	return f.Type == "bool" && f.BoolIndex%8 == 0
}

// DocText returns the documentation lines prefixed with ident.
func (f *Field) DocText(indent string) string {
	return docText(f.Docs, indent)
}

// String returns the qualified name.
func (f *Field) String() string {
	return fmt.Sprintf("%s.%s", f.Struct, f.Name)
}

func docText(docs []string, indent string) string {
	if len(docs) == 0 {
		return ""
	}

	var buf bytes.Buffer
	for i, s := range docs {
		if i != 0 {
			buf.WriteByte('\n')
		}
		if !strings.HasPrefix(s, "// ") {
			continue
		}

		buf.WriteString(indent)
		buf.WriteString(s[3:])
	}

	return buf.String()
}

// StructsByQName maps each Struct to its respective qualified name
// (as in <package>.<type>).
func (p Packages) StructsByQName() map[string]*Struct {
	var n int
	for _, pkg := range p {
		n += len(pkg.Structs)
	}
	m := make(map[string]*Struct, n)

	for _, pkg := range p {
		for _, t := range pkg.Structs {
			qName := t.String()
			if _, ok := m[qName]; ok {
				panic(qName + " dupe")
			}
			m[qName] = t
		}
	}
	return m
}

// FieldsByQName maps each Field to its respective qualified name
// (as in <package>.<type>.<field>).
func (p Packages) FieldsByQName() map[string]*Field {
	var n int
	for _, pkg := range p {
		for _, t := range pkg.Structs {
			n += len(t.Fields)
		}
	}
	m := make(map[string]*Field, n)

	for _, pkg := range p {
		for _, t := range pkg.Structs {
			for _, f := range t.Fields {
				qName := f.String()
				if _, ok := m[qName]; ok {
					panic(qName + " dupe")
				}
				m[qName] = f
			}
		}
	}
	return m
}

// TagAllow defines tag options.
type TagAllow int

const (
	TagNone   TagAllow = iota // not allowed
	TagSingle                 // zero or one
	TagMulti                  // any number
)

type TagOptions struct {
	StructAllow TagAllow
	FieldAllow  TagAllow
}

func (p Packages) ApplyTagFile(path string, options TagOptions) error {
	fields := p.FieldsByQName()
	structs := p.StructsByQName()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	r := bufio.NewReader(file)

	for lineNo := 1; ; lineNo++ {
		line, isPrefix, err := r.ReadLine()
		switch err {
		case nil:
			break
		case io.EOF:
			return nil
		default:
			return err
		}
		if isPrefix {
			return fmt.Errorf("parse %s:%d: line exceeds %d bytes", path, lineNo, r.Size())
		}

		// parse line
		line = bytes.TrimLeftFunc(line, unicode.IsSpace)
		if len(line) == 0 || line[0] == '#' {
			continue // empty or comment
		}
		i := bytes.IndexFunc(line, unicode.IsSpace)
		if i < 0 {
			i = len(line)
		}
		qName := line[:i]
		tag := string(bytes.TrimSpace(line[i:]))
		if tag == "" {
			return fmt.Errorf("parse %s:%d: incomplete declaration %q", path, lineNo, line)
		}

		// match qualifier
		if t := structs[string(qName)]; t != nil {
			switch options.StructAllow {
			case TagNone:
				return fmt.Errorf("apply %s:%d: struct tag (on %s) not supported by target language", path, lineNo, qName)
			case TagSingle:
				if len(t.TagAdd) != 0 {
					return fmt.Errorf("apply %s:%d: %s already tagged [duplicate]", path, lineNo, qName)
				}
			}
			t.TagAdd = append(t.TagAdd, tag)
		} else if f := fields[string(qName)]; f != nil {
			switch options.FieldAllow {
			case TagNone:
				return fmt.Errorf("apply %s:%d: field tag (on %s) not supported by target language", path, lineNo, qName)
			case TagSingle:
				if len(f.TagAdd) != 0 {
					return fmt.Errorf("apply %s:%d: %s already tagged [duplicate]", path, lineNo, qName)
				}
			}
			f.TagAdd = append(f.TagAdd, tag)
		} else {
			return p.qNameNotFound(string(qName), path, lineNo)
		}
	}
}

// QNameNotFound narrows the mismatch down with user-friendly errors.
func (p Packages) qNameNotFound(qName string, path string, lineNo int) error {
	segs := strings.SplitN(qName, ".", 4)
	if len(segs) < 2 || len(segs) > 3 {
		return fmt.Errorf("parse %s:%d: invalid qualifier %q; use <package>'.'<type>('.'<field>)", path, lineNo, qName)
	}

	for _, pkg := range p {
		if !strings.EqualFold(pkg.Name, segs[0]) {
			continue
		}
		if pkg.Name != segs[0] {
			return fmt.Errorf("map %s:%d: package not found; case mismatch with %s?", path, lineNo, pkg.Name)
		}

		for _, t := range pkg.Structs {
			if !strings.EqualFold(t.Name, segs[1]) {
				continue
			}
			if t.Name != segs[1] {
				return fmt.Errorf("map %s:%d: type not found; case mismatch with %s?", path, lineNo, t)
			}

			for _, f := range t.Fields {
				if strings.EqualFold(f.Name, segs[2]) {
					return fmt.Errorf("map %s:%d: field not found; case mismatch with %s?", path, lineNo, f)
				}
			}

			return fmt.Errorf("map %s:%d: field %q not in schema", path, lineNo, segs[0]+"."+segs[1]+"."+segs[2])
		}
		return fmt.Errorf("map %s:%d: type %q not in schema", path, lineNo, segs[0]+"."+segs[1])
	}
	return fmt.Errorf("map %s:%d: package %q not in schema", path, lineNo, segs[0])
}
