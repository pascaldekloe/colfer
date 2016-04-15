package colfer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

// Datatypes hold all supported names.
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
	Name    string
	Structs []*Struct
}

// Struct is data structure definition.
type Struct struct {
	Pkg Package
	// Name is the identification token.
	Name string
	// NameTitle is the identification token in title case.
	NameTitle string
	Fields    []*Field
}

// Field is a Struct member definition.
type Field struct {
	// Index is the Struct.Fields position.
	Index int
	// Name is the identification token.
	Name string
	// NameTitle is the identification token in title case.
	NameTitle string
	// Type is the datatype.
	Type string
}

// ReadDefs parses the Colfer files.
func ReadDefs(files []string) ([]*Struct, error) {
	var structs []*Struct

	fileSet := token.NewFileSet()
	for _, file := range files {
		file, err := parser.ParseFile(fileSet, file, nil, 0)
		if err != nil {
			return nil, err
		}

		pkg := Package{}
		pkg.Name = file.Name.Name

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

				o, err := mapStruct(s)
				if err != nil {
					return nil, err
				}
				o.Pkg = pkg
				structs = append(structs, o)
			}
		}
	}

	return structs, nil
}

func mapStruct(src *ast.TypeSpec) (*Struct, error) {
	dst := &Struct{
		Name:      src.Name.Name,
		NameTitle: strings.Title(src.Name.Name),
	}

	s, ok := src.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("colfer: unsupported type %s", reflect.TypeOf(s))
	}

	for i, f := range s.Fields.List {
		field := Field{Index: i}
		dst.Fields = append(dst.Fields, &field)

		if len(f.Names) == 0 {
			return nil, fmt.Errorf("colfer: missing name for field %d", i)
		}
		field.Name = f.Names[0].Name
		field.NameTitle = strings.Title(f.Names[0].Name)

		t, ok := f.Type.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("colfer: unknow type in stuct %q field %d %q: %#v", dst.Name, field.Index, field.Name, field.Type)
		}
		field.Type = t.Name
		if _, ok := Datatypes[field.Type]; !ok {
			return nil, fmt.Errorf("colfer: unknown datatype %q in struct %q field %d %q", field.Type, dst.Name, field.Index, field.Name)
		}
	}

	return dst, nil
}
