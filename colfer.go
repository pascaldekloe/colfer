// Package colfer provides the schema interpretation.
package colfer

//go:generate go run ./cmd/colf/main.go -b rpc go rpc/header.colf

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// datatypes holds all supported names.
var datatypes = map[string]struct{}{
	"bool":      {},
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

// Package is a named definition bundle.
type Package struct {
	// Name is the identification token.
	Name string
	// NameNative is the language specific identification token
	NameNative string
	Structs    []*Struct
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
	Name   string
	Fields []*Field
}

// NameTitle returns the identification token in title case.
func (s *Struct) NameTitle() string {
	return strings.Title(s.Name)
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
	// Type is the datatype.
	Type string
	// TypeNative is the language specific datatype placeholder.
	TypeNative string
	// TypeRef is the Colfer data structure reference.
	TypeRef *Struct
	// TypeList flags whether the datatype is a list.
	TypeList bool
}

// NameTitle gets the identification token in title case.
func (f *Field) NameTitle() string {
	return strings.Title(f.Name)
}

// String returns the qualified name.
func (f *Field) String() string {
	return fmt.Sprintf("%s.%s", f.Struct, f.Name)
}

// ReadDefs parses the Colfer files.
func ReadDefs(files []string) ([]*Package, error) {
	var packages []*Package

	fileSet := token.NewFileSet()
	for _, file := range files {
		file, err := parser.ParseFile(fileSet, file, nil, 0)
		if err != nil {
			return nil, err
		}

		var pkg *Package
		for _, p := range packages {
			if p.Name == file.Name.Name {
				pkg = p
				break
			}
		}
		if pkg == nil {
			pkg = &Package{Name: file.Name.Name}
			packages = append(packages, pkg)
		}

		for _, decl := range file.Decls {
			d, ok := decl.(*ast.GenDecl)
			if !ok {
				return nil, fmt.Errorf("colfer: unsupported declaration type %s", reflect.TypeOf(decl))
			}

			for _, spec := range d.Specs {
				s, ok := spec.(*ast.TypeSpec)
				if !ok {
					return nil, fmt.Errorf("colfer: unsupported specification type %s", reflect.TypeOf(spec))
				}

				if err := addStruct(pkg, s); err != nil {
					return nil, err
				}
			}
		}
	}

	names := make(map[string]*Struct)
	for _, pkg := range packages {
		for _, s := range pkg.Structs {
			qname := s.String()
			if _, ok := names[qname]; ok {
				return nil, fmt.Errorf("colfer: duplicate struct definition %q", qname)
			}
			names[qname] = s
		}
	}

	for _, pkg := range packages {
		for _, s := range pkg.Structs {
			for _, f := range s.Fields {
				t := f.Type
				_, ok := datatypes[t]
				if ok {
					if f.TypeList && t != "text" {
						return nil, fmt.Errorf("colfer: unsupported lists type %s for field %q", t, f)
					}
					continue
				}
				if f.TypeRef, ok = names[t]; ok {
					continue
				}
				if f.TypeRef, ok = names[pkg.Name+"."+t]; ok {
					continue
				}
				return nil, fmt.Errorf("colfer: unknown datatype for field %q", f)
			}
		}
	}

	return packages, nil
}

func addStruct(pkg *Package, src *ast.TypeSpec) error {
	dst := &Struct{
		Pkg:  pkg,
		Name: src.Name.Name,
	}
	pkg.Structs = append(pkg.Structs, dst)

	s, ok := src.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("colfer: unsupported type %s", reflect.TypeOf(s))
	}

	for i, f := range s.Fields.List {
		field := Field{Struct: dst, Index: i}
		dst.Fields = append(dst.Fields, &field)

		if len(f.Names) == 0 {
			return fmt.Errorf("colfer: missing name for field %d", i)
		}
		field.Name = f.Names[0].Name

		t := f.Type
		if at, ok := t.(*ast.ArrayType); ok {
			t = at.Elt
			field.TypeList = true
		}

		id, ok := t.(*ast.Ident)
		if !ok {
			return fmt.Errorf("colfer: unknown datatype declaration for field %q: %#v", field, t)
		}
		field.Type = id.Name
	}

	return nil
}
