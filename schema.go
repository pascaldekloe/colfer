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

	names := make(map[string]*Struct)
	for _, pkg := range packages {
		for _, s := range pkg.Structs {
			qname := s.String()
			if dupe, ok := names[qname]; ok {
				return nil, fmt.Errorf("colfer: duplicate struct definition %q in file %s and %s", qname, dupe.SchemaFile, s.SchemaFile)
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
					if f.TypeList {
						switch t {
						case "int32", "int64":
							fmt.Println("colfer: WARNING: integer lists are Go only at the moment")
						case "float32", "float64", "text", "binary":
						default:
							return nil, fmt.Errorf("colfer: unsupported lists type %q for field %s", t, f)
						}
					}
					continue
				}
				if f.TypeRef, ok = names[t]; ok {
					continue
				}
				if f.TypeRef, ok = names[pkg.Name+"."+t]; ok {
					continue
				}
				return nil, fmt.Errorf("colfer: unknown datatype %q for field %s", t, f)
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
		switch t := spec.Type.(type) {
		default:
			return fmt.Errorf("colfer: unsupported data type %T", t)
		case *ast.StructType:
			s := &Struct{Pkg: pkg, Name: spec.Name.Name, SchemaFile: path.Base(schemaPath)}
			pkg.Structs = append(pkg.Structs, s)

			s.Docs = append(docs(decl.Doc), docs(spec.Doc)...)
			if err := mapStruct(s, t); err != nil {
				return err
			}
		}
	}

	return nil
}

func mapStruct(dst *Struct, src *ast.StructType) error {
	for i, f := range src.Fields.List {
		field := &Field{Struct: dst, Index: i}
		dst.Fields = append(dst.Fields, field)

		if len(f.Names) == 0 {
			return fmt.Errorf("colfer: field %d from %s has no name", i, dst)
		}
		field.Name = f.Names[0].Name

		if f.Tag != nil {
			return fmt.Errorf("colfer: illegal tag %s on %s", f.Tag.Value, field)
		}

		field.Docs = docs(f.Doc)

		expr := f.Type
		for {
			switch t := expr.(type) {
			case *ast.ArrayType:
				expr = t.Elt
				field.TypeList = true
				continue
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
			break
		}
	}

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
