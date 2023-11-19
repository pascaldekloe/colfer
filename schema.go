package colfer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"strconv"
)

// FormatFile normalizes the structure.
// The content is expected to be syntactically correct.
func FormatFile(path string) (changed bool, err error) {
	orig, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	clean, err := format.Source(orig)
	if err != nil {
		return false, fmt.Errorf("colfer: format %q: %s", path, err)
	}

	if bytes.Equal(orig, clean) {
		return false, nil
	}

	if err := ioutil.WriteFile(path, clean, 0644); err != nil {
		return false, err
	}
	return true, nil
}

// ParseFiles returns the schema definitions.
func ParseFiles(paths ...string) (Packages, error) {
	var packages Packages

	fileSet := token.NewFileSet()
	for _, schemaPath := range paths {
		fileAST, err := parser.ParseFile(fileSet, schemaPath, nil, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return nil, err
		}

		var pkg *Package
		for _, p := range packages {
			if p.Name == fileAST.Name.Name {
				pkg = p
				break
			}
		}
		if pkg == nil {
			pkg = &Package{Name: fileAST.Name.Name}
			packages = append(packages, pkg)
		}

		pkg.SchemaFiles = append(pkg.SchemaFiles, path.Base(schemaPath))

		pkg.Docs = append(pkg.Docs, docs(fileAST.Doc)...)

		// switch through the AST types
		for _, decl := range fileAST.Decls {
			switch decl := decl.(type) {
			default:
				return nil, fmt.Errorf("colfer: unsupported declaration type %T", decl)
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					if err := addSpec(pkg, decl, spec, schemaPath); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	structPerName := make(map[string]*Struct)
	for _, pkg := range packages {
		for _, t := range pkg.Structs {
			qname := t.String()
			if dupe, ok := structPerName[qname]; ok {
				return nil, fmt.Errorf("colfer: duplicate struct definition %q in file %s and %s", qname, dupe.SchemaFile, t.SchemaFile)
			}
			structPerName[qname] = t
		}
	}

	// type checks
	for _, pkg := range packages {
		for _, t := range pkg.Structs {
			for _, f := range t.Fields {
				if _, ok := datatypes[f.Type]; ok {
					continue // pass
				}

				ref, ok := structPerName[f.Type]
				if !ok {
					ref, ok = structPerName[pkg.Name+"."+f.Type]
				}
				if !ok {
					return nil, fmt.Errorf("colfer: unknown datatype %q on field %s", f.Type, f)
				}
				f.TypeRef = ref
			}
		}
	}

	return packages, nil
}

func addSpec(pkg *Package, decl *ast.GenDecl, spec ast.Spec, schemaPath string) error {
	switch spec := spec.(type) {
	default:
		return fmt.Errorf("colfer: unsupported specification type %T", spec)
	case *ast.TypeSpec:
		switch specType := spec.Type.(type) {
		default:
			return fmt.Errorf("colfer: unsupported data type %T", specType)
		case *ast.StructType:
			t := &Struct{Pkg: pkg, Name: spec.Name.Name, SchemaFile: path.Base(schemaPath)}
			for _, pt := range pkg.Structs {
				if pt.Name == t.Name {
					return fmt.Errorf("colfer: duplicate %s declaration", pt)
				}
			}
			pkg.Structs = append(pkg.Structs, t)

			t.Docs = append(docs(decl.Doc), docs(spec.Doc)...)
			if err := mapStruct(t, specType); err != nil {
				return err
			}
		}
	}

	return nil
}

func mapStruct(dst *Struct, src *ast.StructType) error {
	if len(src.Fields.List) == 0 {
		return fmt.Errorf("colfer: %s has no fields", dst)
	}

	var colferFieldIndex int // counts array elements individually
	for i, f := range src.Fields.List {
		field := &Field{Struct: dst, Index: colferFieldIndex}
		dst.Fields = append(dst.Fields, field)
		colferFieldIndex++

		if len(f.Names) == 0 {
			return fmt.Errorf("colfer: field %d from %s has no name", i, dst)
		}
		field.Name = f.Names[0].Name

		if f.Tag != nil {
			return fmt.Errorf("colfer: illegal tag %s on %s", f.Tag.Value, field)
		}

		field.Docs = docs(f.Doc)

		ftype := f.Type
		if array, ok := ftype.(*ast.ArrayType); ok {
			ftype = array.Elt

			if array.Len != nil {
				l, ok := array.Len.(*ast.BasicLit)
				if !ok {
					return fmt.Errorf("colfer: unknown array lenth type %T for field %s", array.Len, field)
				}
				n, err := strconv.Atoi(l.Value)
				if err != nil || n < 2 || n > 256 {
					return fmt.Errorf("colfer: array size %q for field %s not within range [2, 256]", l.Value, field)
				}

				field.ElementCount = n
				colferFieldIndex += n - 1
			} else {
				field.TypeList = true
			}
		}

		switch t := ftype.(type) {
		case *ast.Ident:
			field.Type = t.Name

		case *ast.SelectorExpr:
			switch pkgIdent := t.X.(type) {
			case *ast.Ident:
				field.Type = pkgIdent.Name + "." + t.Sel.Name
			default:
				return fmt.Errorf("colfer: unknown datatype selector expression %T for field %s", pkgIdent, field)
			}

		default:
			return fmt.Errorf("colfer: unknown datatype declaration %T for field %s", t, field)
		}
	}

	dst.SetFixedPositions()

	return nil
}

func docs(g *ast.CommentGroup) []string {
	var a []string
	if g != nil {
		for _, c := range g.List {
			a = append(a, c.Text)
		}
	}
	return a
}
