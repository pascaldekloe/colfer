package colfer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// Datatypes holds all supported names.
var Datatypes = map[string]struct{}{
	"bool":      struct{}{},
	"uint32":    struct{}{},
	"uint64":    struct{}{},
	"int32":     struct{}{},
	"int64":     struct{}{},
	"float32":   struct{}{},
	"float64":   struct{}{},
	"timestamp": struct{}{},
	"text":      struct{}{},
	"binary":    struct{}{},
}

// Package is a named definition bundle.
type Package struct {
	// Name is the identification token.
	Name string
	// NameNative is the language specific identification token
	NameNative string
	Structs    []*Struct
}

// Struct is a data structure definition.
type Struct struct {
	Pkg *Package
	// Name is the identification token.
	Name   string
	Fields []*Field
}

// NameTitle gets the identification token in title case.
func (s *Struct) NameTitle() string {
	return strings.Title(s.Name)
}

func (s *Struct) String() string {
	return fmt.Sprintf("%s.%s", s.Pkg.Name, s.Name)
}

// Field is a Struct member definition.
type Field struct {
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
}

// NameTitle gets the identification token in title case.
func (f *Field) NameTitle() string {
	return strings.Title(f.Name)
}

func (f *Field) String() string {
	return fmt.Sprintf("%s:%s", f.Name, f.Type)
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

	if err := linkStructs(packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func linkStructs(packages []*Package) error {
	names := make(map[string]*Struct)

	for _, pkg := range packages {
		for _, s := range pkg.Structs {
			qname := s.String()
			if _, ok := names[qname]; ok {
				return fmt.Errorf("colfer: duplicate struct definition %q", qname)
			}
			names[qname] = s
		}
	}

	for _, pkg := range packages {
		for _, s := range pkg.Structs {
			for _, f := range s.Fields {
				t := f.Type
				_, ok := Datatypes[t]
				if ok {
					continue
				}
				if f.TypeRef, ok = names[t]; ok {
					continue
				}
				if f.TypeRef, ok = names[pkg.Name+"."+t]; ok {
					continue
				}
				return fmt.Errorf("colfer: unknown datatype in struct %q field %q", s, f)
			}
		}
	}

	return nil
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
		field := Field{Index: i}
		dst.Fields = append(dst.Fields, &field)

		if len(f.Names) == 0 {
			return fmt.Errorf("colfer: missing name for field %d", i)
		}
		field.Name = f.Names[0].Name

		t, ok := f.Type.(*ast.Ident)
		if !ok {
			return fmt.Errorf("colfer: unknow type in stuct %q field %d %q: %#v", dst, field.Index, field.Name, field.Type)
		}
		field.Type = t.Name
	}

	return nil
}
