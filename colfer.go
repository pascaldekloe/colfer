// Package colfer provides schema definitions.
package colfer

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// datatypes holds all supported names.
var datatypes = map[string]struct{}{
	"bool":      {},
	"uint8":     {},
	"uint16":    {},
	"uint32":    {},
	"uint64":    {},
	"int32":     {},
	"int64":     {},
	"float32":   {},
	"float64":   {},
	"timestamp": {},
	"text":      {},
	"binary":    {},
}

type packages []*Package

func (p packages) Len() int           { return len(p) }
func (p packages) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p packages) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Package is a named definition bundle.
type Package struct {
	// Name is the identification token.
	Name string
	// NameNative is the language specific identification token
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Structs are the type definitions.
	Structs []*Struct
	// SchemaFiles are the source filenames.
	SchemaFiles []string
	// SizeMax is the uper limit expression.
	SizeMax string
	// ListMax is the uper limit expression.
	ListMax string
	// SuperClass is the superclass.
	SuperClass string
	// SuperClassNative is the language specific superclass.
	SuperClassNative string
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
func (p *Package) Refs() []*Package {
	found := make(map[*Package]struct{})
	for _, s := range p.Structs {
		for _, f := range s.Fields {
			if f.TypeRef != nil && f.TypeRef.Pkg != p {
				found[f.TypeRef.Pkg] = struct{}{}
			}
		}
	}

	var refs packages
	for r := range found {
		refs = append(refs, r)
	}
	sort.Sort(refs)
	return refs
}

// HasFloat returns whether p has one or more floating point fields.
func (p *Package) HasFloat() bool {
	for _, s := range p.Structs {
		if s.HasFloat() {
			return true
		}
	}
	return false
}

// HasTimestamp returns whether p has one or more timestamp fields.
func (p *Package) HasTimestamp() bool {
	for _, s := range p.Structs {
		if s.HasTimestamp() {
			return true
		}
	}
	return false
}

// HasList returns whether p has one or more list fields.
func (p *Package) HasList() bool {
	for _, s := range p.Structs {
		if s.HasList() {
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
	// NameNative is the language specific identification token
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Fields are the elements in order of appearance.
	Fields []*Field
	// SchemaFile is the source filename.
	SchemaFile string
}

// NameTitle returns the identification token in title case.
func (s *Struct) NameTitle() string {
	return strings.Title(s.Name)
}

// DocText returns the documentation lines prefixed with ident.
func (s *Struct) DocText(indent string) string {
	return docText(s.Docs, indent)
}

// String returns the qualified name.
func (s *Struct) String() string {
	return fmt.Sprintf("%s.%s", s.Pkg.Name, s.Name)
}

// HasFloat returns whether s has one or more floating point fields.
func (s *Struct) HasFloat() bool {
	for _, f := range s.Fields {
		if f.Type == "float32" || f.Type == "float64" {
			return true
		}
	}
	return false
}

// HasText returns whether s has one or more text fields.
func (s *Struct) HasText() bool {
	for _, f := range s.Fields {
		if f.Type == "text" {
			return true
		}
	}
	return false
}

// HasBinary returns whether s has one or more binary fields.
func (s *Struct) HasBinary() bool {
	for _, f := range s.Fields {
		if f.Type == "binary" {
			return true
		}
	}
	return false
}

// HasBinaryList returns whether s has one or more binary list fields.
func (s *Struct) HasBinaryList() bool {
	for _, f := range s.Fields {
		if f.Type == "binary" && f.TypeList {
			return true
		}
	}
	return false
}

// HasTimestamp returns whether s has one or more timestamp fields.
func (s *Struct) HasTimestamp() bool {
	for _, f := range s.Fields {
		if f.Type == "timestamp" {
			return true
		}
	}
	return false
}

// HasList returns whether s has one or more list fields.
func (s *Struct) HasList() bool {
	for _, f := range s.Fields {
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
	// NameNative is the language specific identification token
	NameNative string
	// Docs are the documentation texts.
	Docs []string
	// Type is the datatype.
	Type string
	// TypeNative is the language specific datatype placeholder.
	TypeNative string
	// TypeRef is the Colfer data structure reference.
	TypeRef *Struct
	// TypeList flags whether the datatype is a list.
	TypeList bool
}

// NameTitle returns the identification token in title case.
func (f *Field) NameTitle() string {
	return strings.Title(f.Name)
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
	for _, s := range docs {
		if !strings.HasPrefix(s, "//") {
			continue
		}

		buf.WriteString(indent)
		buf.WriteString(s[2:])
		buf.WriteByte('\n')
	}

	return buf.String()
}
