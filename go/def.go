package colfer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"
)

// Object is a data structure definition.
type Object struct {
	Package string
	Name    string
	Fields  []*Field
}

// Field is an Object item definition.
type Field struct {
	Num  int
	Name string
	Type string
}

// ReadDefs parses the Colfer files.
func ReadDefs() ([]*Object, error) {
	var objects []*Object
	fileSet := token.NewFileSet()

	colfFiles, err := filepath.Glob("*.colf")
	if err != nil {
		return nil, err
	}
	for _, path := range colfFiles {
		file, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return nil, err
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

				o, err := mapObject(s)
				if err != nil {
					return nil, err
				}
				objects = append(objects, o)

				o.Package = file.Name.Name
			}
		}
	}

	return objects, nil
}

func mapObject(src *ast.TypeSpec) (*Object, error) {
	dst := &Object{
		Name: strings.Title(src.Name.Name),
	}

	s, ok := src.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("colfer: unsupported type %s", reflect.TypeOf(s))
	}
	for i, f := range s.Fields.List {
		t, ok := f.Type.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("colfer: unknow type for field %d: %#v", i, f.Type)
		}
		if len(f.Names) == 0 {
			return nil, fmt.Errorf("colfer: missing name for field %d", i)
		}

		dst.Fields = append(dst.Fields, &Field{
			Num:  i,
			Name: strings.Title(f.Names[0].Name),
			Type: t.Name,
		})
	}

	return dst, nil
}
