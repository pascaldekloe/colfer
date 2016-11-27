// Package colfer provides the schema interpretation.
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
	// name is the identification token.
	name string
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
	return strings.Title(f.name)
}

// DocText returns the documentation lines prefixed with ident.
func (f *Field) DocText(indent string) string {
	return docText(f.Docs, indent)
}

// String returns the qualified name.
func (f *Field) String() string {
	return fmt.Sprintf("%s.%s", f.Struct, f.name)
}

// Format normalizes the file's content.
// The content of file is expected to be syntactically correct.
func Format(file string) (changed bool, err error) {
	orig, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}

	clean, err := format.Source(orig)
	if err != nil {
		return false, fmt.Errorf("colfer: format %q: %s", file, err)
	}

	if bytes.Equal(orig, clean) {
		return false, nil
	}

	if err := ioutil.WriteFile(file, clean, 0644); err != nil {
		return false, err
	}
	return true, nil
}

// ReadDefs parses schema files.
func ReadDefs(files []string) ([]*Package, error) {
	var packages []*Package

	fileSet := token.NewFileSet()
	for _, file := range files {
		fileAST, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments|parser.AllErrors)
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

		pkg.SchemaFiles = append(pkg.SchemaFiles, path.Base(file))

		pkg.Docs = append(pkg.Docs, docs(fileAST.Doc)...)

		// switch through the AST types
		for _, decl := range fileAST.Decls {
			switch decl := decl.(type) {
			default:
				return nil, fmt.Errorf("colfer: unsupported declaration type %T", decl)
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					if err := addSpec(pkg, decl, spec, file); err != nil {
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
					if f.TypeList && t != "text" && t != "binary" {
						return nil, fmt.Errorf("colfer: unsupported lists type %q for field %s", t, f.String())
					}
					continue
				}
				if f.TypeRef, ok = names[t]; ok {
					continue
				}
				if f.TypeRef, ok = names[pkg.Name+"."+t]; ok {
					continue
				}
				return nil, fmt.Errorf("colfer: unknown datatype %q for field %s", t, f.String())
			}
		}
	}

	return packages, nil
}

func addSpec(pkg *Package, decl *ast.GenDecl, spec ast.Spec, file string) error {
	switch spec := spec.(type) {
	default:
		return fmt.Errorf("colfer: unsupported specification type %T", spec)
	case *ast.TypeSpec:
		switch t := spec.Type.(type) {
		default:
			return fmt.Errorf("colfer: unsupported data type %T", t)
		case *ast.StructType:
			s := &Struct{Pkg: pkg, Name: spec.Name.Name, SchemaFile: path.Base(file)}
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
		field := Field{Struct: dst, Index: i}
		dst.Fields = append(dst.Fields, &field)

		if len(f.Names) == 0 {
			return fmt.Errorf("colfer: missing name for field %d", i)
		}
		field.name = f.Names[0].Name

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
					return fmt.Errorf("colfer: unknown datatype selector expression %T for field %s", pkgIdent, field.String())
				}
			default:
				return fmt.Errorf("colfer: unknown datatype declaration %T for field %s", t, field.String())
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
