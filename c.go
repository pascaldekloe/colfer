package colfer

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pascaldekloe/name"
)

// CKeywords are the reserved tokens for C code.
// Some entries are redundant due to the use of a Go parser.
var cKeywords = map[string]struct{}{
	"auto": {}, "break": {}, "case": {}, "char": {},
	"const": {}, "continue": {}, "default": {}, "do": {},
	"double": {}, "else": {}, "enum": {}, "extern": {},
	"float": {}, "for": {}, "goto": {}, "if": {},
	"int": {}, "long": {}, "register": {}, "return": {},
	"short": {}, "signed": {}, "sizeof": {}, "static": {},
	"struct": {}, "switch": {}, "typedef": {}, "union": {},
	"unsigned": {}, "void": {}, "volatile": {}, "while": {},
}

// GenerateC writes the code into file "Colfer.h" and "Colfer.c".
func GenerateC(basedir string, packages Packages) error {
	for _, p := range packages {
		for _, t := range p.Structs {
			t.NameNative = strings.ToLower(name.SnakeCase(p.Name + "_" + t.Name))

			for _, f := range t.Fields {
				f.NameNative = strings.ToLower(name.SnakeCase(f.Name))
				if _, ok := cKeywords[f.NameNative]; ok {
					f.NameNative += "_"
				}

				switch f.Type {
				case "bool":
					f.TypeNative = "char"
				case "uint8", "uint16", "uint32", "uint64", "int32", "int64":
					f.TypeNative = f.Type + "_t"
				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"
				case "timestamp":
					f.TypeNative = "timespec"
				case "binary", "text":
					f.TypeNative = "colfer_" + f.Type
				}
			}
		}
	}

	if err := os.MkdirAll(basedir, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(basedir, "Colfer.h"))
	if err != nil {
		return err
	}
	if err := template.Must(template.New("C-header").Parse(cHeaderTemplate)).Execute(f, packages); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(basedir, "Colfer.c"))
	if err != nil {
		return err
	}
	if err := template.Must(template.New("C").Parse(cTemplate)).Execute(f, packages); err != nil {
		return err
	}
	return f.Close()
}

//go:embed template/c-header.txt
var cHeaderTemplate string

//go:embed template/c.txt
var cTemplate string
