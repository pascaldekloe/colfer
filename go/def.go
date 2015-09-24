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

// Package is a bundle of one or more Object definitions.
type Package struct {
	Name    string
	Objects []*Object
}

// Object is a data structure definition.
type Object struct {
	Name   string
	Fields []*Field
}

// Field is an Object item definition.
type Field struct {
	Num  int
	Name string
	Type string
}

// ReadDefs parses the Colfer files.
func ReadDefs() (*Package, error) {
	pkg := new(Package)
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
		if pkgName := file.Name.Name; pkg.Name == "" {
			pkg.Name = pkgName
		} else if pkgName != pkg.Name {
			return nil, fmt.Errorf("colfer: package mismatch: %q and %q", pkgName, pkg.Name)
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
				pkg.Objects = append(pkg.Objects, o)
			}
		}
	}

	return pkg, nil
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
