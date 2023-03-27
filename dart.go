package colfer

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/pascaldekloe/name"
)

// dartKeywords are the reserved tokens for Dart code.
// Some entries are redundant due to the use of a Go parser.
var dartKeywords = map[string]struct{}{
	"abstract": {}, "as": {}, "assert": {}, "async": {}, "await": {},
	"break": {}, "case": {}, "catch": {}, "class": {}, "const": {},
	"continue": {}, "covariant": {}, "default": {}, "deferred": {},
	"do": {}, "dynamic": {}, "else": {}, "enum": {}, "export": {},
	"extends": {}, "extension": {}, "external": {}, "factory": {},
	"false": {}, "final": {}, "finally": {}, "for": {}, "Function": {},
	"get": {}, "hide": {}, "if": {}, "implements": {}, "import": {},
	"in": {}, "interface": {}, "is": {}, "library": {}, "mixin": {},
	"new": {}, "null": {}, "on": {}, "operator": {}, "part": {},
	"rethrow": {}, "return": {}, "set": {}, "show": {}, "static": {},
	"super": {}, "switch": {}, "sync": {}, "this": {}, "throw": {},
	"true": {}, "try": {}, "typedef": {}, "var": {}, "void": {},
	"while": {}, "with": {}, "yield": {},
	"other": {}, "marshalTo": {}, "marshalLen": {}, "unmarshal": {},
}

// GenerateDart writes the code into file "Colfer.dart".
func GenerateDart(basedir string, packages Packages) error {
	if err := os.MkdirAll(basedir, os.ModeDir|os.ModePerm); err != nil {
		return err
	}

	t := template.New("dart-code")
	template.Must(t.Parse(dartCode))
	template.Must(t.New("marshal").Parse(dartMarshal))
	template.Must(t.New("marshal-len").Parse(dartMarshalLen))
	template.Must(t.New("unmarshal").Parse(dartUnmarshal))

	nativeTypes := map[string]string{
		"text":      "String",
		"binary":    "Uint8List",
		"timestamp": "DateTime?",
		"uint8":     "int",
		"uint16":    "int",
		"uint32":    "int",
		"uint64":    "int",
		"int32":     "int",
		"int64":     "int",
		"float32":   "double",
		"float64":   "double",
	}
	nativeListTypes := map[string]string{
		"uint8":   "Uint8List",
		"uint16":  "Uint16List",
		"uint32":  "Uint32List",
		"uint64":  "Uint64List",
		"int32":   "Int32List",
		"int64":   "Int64List",
		"float32": "Float32List",
		"float64": "Float64List",
	}

	for _, p := range packages {
		p.NameNative = p.Name
		if _, ok := dartKeywords[p.NameNative]; ok {
			p.NameNative += "_0"
		}
		for _, t := range p.Structs {
			t.NameNative = name.CamelCase(t.Name, true)
			if _, ok := dartKeywords[t.NameNative]; ok {
				t.NameNative += "_0"
			}
			for _, f := range t.Fields {
				f.NameNative = name.CamelCase(f.Name, false)
				if _, ok := dartKeywords[f.NameNative]; ok {
					f.NameNative += "_0"
				}
			}
		}
	}

	for _, p := range packages {
		for _, t := range p.Structs {
			for _, f := range t.Fields {
				if f.TypeRef != nil {
					f.TypeNative = f.TypeRef.NameNative
					if f.TypeRef.Pkg != p {
						f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
					}
					continue
				}
				if f.TypeList {
					if nativeType, ok := nativeListTypes[f.Type]; ok {
						f.TypeNative = nativeType
					} else {
						f.TypeNative = nativeTypes[f.Type]
					}
				} else {
					if nativeType, ok := nativeTypes[f.Type]; ok {
						f.TypeNative = nativeType
					} else {
						f.TypeNative = f.Type
					}
				}
			}
		}
		path := filepath.Join(basedir, p.Name)
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(path, "Colfer.dart"))
		if err != nil {
			return err
		}
		defer f.Close()
		if err = t.Execute(f, p); err != nil {
			return err
		}
		if err = f.Sync(); err != nil {
			return err
		}
	}

	res, err := exec.Command("dart", "format", basedir, "--fix", "-l", "100").Output()
	fmt.Printf("%s", res)
	return err
}


//go:embed template/dart.txt
var dartCode string

//go:embed template/dart-marshal.txt
var dartMarshal string

//go:embed template/dart-marshal-len.txt
var dartMarshalLen string

//go:embed template/dart-unmarshal.txt
var dartUnmarshal string
