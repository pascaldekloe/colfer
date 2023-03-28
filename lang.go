package colfer

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/pascaldekloe/name"
	"golang.org/x/mod/modfile"
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
		"bool":      "bool",
		"uint8":     "int",
		"uint16":    "int",
		"uint32":    "int",
		"uint64":    "int",
		"int32":     "int",
		"int64":     "int",
		"float32":   "double",
		"float64":   "double",
		"timestamp": "DateTime",
		"text":      "String",
		"binary":    "Uint8List",
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
				switch {
				case f.TypeRef != nil:
					f.TypeNative = f.TypeRef.NameNative
					if f.TypeRef.Pkg != p {
						f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
					}
				case f.TypeList:
					if nativeType, ok := nativeListTypes[f.Type]; ok {
						f.TypeNative = nativeType
					} else {
						f.TypeNative = "List<" + nativeTypes[f.Type] + ">"
					}
				default:
					f.TypeNative = nativeTypes[f.Type]
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

// ECMAKeywords are the reserved tokens for ECMA Script.
// Some entries are redundant due to the use of a Go parser.
var eCMAKeywords = map[string]struct{}{
	"break": {}, "case": {}, "catch": {}, "class": {},
	"const": {}, "continue": {}, "debugger": {}, "default": {},
	"delete": {}, "do": {}, "else": {}, "enum": {},
	"export": {}, "extends": {}, "finally": {}, "for": {},
	"function": {}, "if": {}, "import": {}, "in": {},
	"instanceof": {}, "new": {}, "return": {}, "super": {},
	"switch": {}, "this": {}, "throw": {}, "try": {},
	"typeof": {}, "var": {}, "void": {}, "while": {},
	"with": {}, "yield": {},
}

// GenerateECMA writes the code into file "Colfer.js".
func GenerateECMA(basedir string, packages Packages) error {
	for _, p := range packages {
		p.NameNative = strings.Replace(p.Name, "/", "_", -1)
		if _, ok := eCMAKeywords[p.NameNative]; ok {
			p.NameNative += "_"
		}

		for _, t := range p.Structs {
			t.NameNative = name.CamelCase(t.Name, true)
			for _, f := range t.Fields {
				f.NameNative = name.CamelCase(f.Name, false)
				if _, ok := eCMAKeywords[f.NameNative]; ok {
					f.NameNative += "_"
				}
			}
		}
	}

	t := template.New("ecma-code")
	template.Must(t.Parse(ecmaCode))
	template.Must(t.New("marshal").Parse(ecmaMarshal))
	template.Must(t.New("unmarshal").Parse(ecmaUnmarshal))

	if err := os.MkdirAll(basedir, os.ModeDir|os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(basedir, "Colfer.js"))
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, packages)
}

//go:embed template/ecma.txt
var ecmaCode string

//go:embed template/ecma-marshal.txt
var ecmaMarshal string

//go:embed template/ecma-unmarshal.txt
var ecmaUnmarshal string

// GoMod looks for a Go modules definition.
// ModDir is the root directory when found.
// ModPkg is the root package when found.
func goMod(dir string) (modDir, modPkg string, err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}

	for n := 0; n < 32; n++ {
		path := filepath.Join(dir, "go.mod")
		text, err := ioutil.ReadFile(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return "", "", err
			}

			// The path does not end in a separator
			// unless it is the root directory.
			if dir[len(dir)-1] == filepath.Separator {
				break
			}
			// try parent directory
			dir = filepath.Dir(dir)
			continue
		}

		return dir, modfile.ModulePath(text), nil
	}

	return "", "", nil // not found
}

// GenerateGo writes the code into file "Colfer.go".
func GenerateGo(basedir string, packages Packages) error {
	t := template.New("go-code")
	template.Must(t.Parse(goCode))
	template.Must(t.New("marshal-field").Parse(goMarshalField))
	template.Must(t.New("marshal-field-len").Parse(goMarshalFieldLen))
	template.Must(t.New("unmarshal-field").Parse(goUnmarshalField))
	template.Must(t.New("unmarshal-varint").Parse(goUnmarshalVarint))

	modDir, modPkg, err := goMod(basedir)
	if err != nil {
		return err
	}

	for _, p := range packages {
		p.NameNative = p.Name[strings.LastIndexByte(p.Name, '/')+1:]
		for _, t := range p.Structs {
			t.NameNative = name.CamelCase(t.Name, true)
			for _, f := range t.Fields {
				f.NameNative = name.CamelCase(f.Name, true)
			}
		}
	}

	for _, p := range packages {
		for _, t := range p.Structs {
			for _, f := range t.Fields {
				switch f.Type {
				default:
					if f.TypeRef == nil {
						f.TypeNative = f.Type
					} else {
						f.TypeNative = f.TypeRef.NameNative
						if f.TypeRef.Pkg != p {
							f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
						}
					}
				case "timestamp":
					f.TypeNative = "time.Time"
				case "text":
					f.TypeNative = "string"
				case "binary":
					f.TypeNative = "[]byte"
				}
			}
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, p); err != nil {
			return err
		}

		path := filepath.Join(basedir, p.Name)
		if modPkg != "" && strings.HasPrefix(p.Name, modPkg+"/") {
			path = filepath.Join(modDir, p.Name[len(modPkg):])
		}
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}

		path = filepath.Join(path, "Colfer.go")
		if err := ioutil.WriteFile(path, buf.Bytes(), 0666); err != nil {
			return err
		}

		if _, err := FormatFile(path); err != nil {
			return err
		}
	}
	return nil
}

//go:embed template/go.txt
var goCode string

//go:embed template/go-marshal-field.txt
var goMarshalField string

//go:embed template/go-marshal-field-len.txt
var goMarshalFieldLen string

//go:embed template/go-unmarshal-field.txt
var goUnmarshalField string

//go:embed template/go-unmarshal-varint.txt
var goUnmarshalVarint string

// JavaKeywords are the reserved tokens for Java code.
// Some entries are redundant due to the use of a Go parser.
var javaKeywords = map[string]struct{}{
	"abstract": {}, "assert": {}, "boolean": {}, "break": {},
	"byte": {}, "case": {}, "catch": {}, "char": {},
	"class": {}, "const": {}, "continue": {}, "default": {},
	"do": {}, "double": {}, "else": {}, "enum": {},
	"extends": {}, "final": {}, "finally": {}, "float": {},
	"for": {}, "goto": {}, "if": {}, "implements": {},
	"import": {}, "instanceof": {}, "int": {}, "interface": {},
	"long": {}, "native": {}, "new": {}, "package": {},
	"private": {}, "protected": {}, "public": {}, "return": {},
	"short": {}, "static": {}, "strictfp": {}, "super": {},
	"switch": {}, "synchronized": {}, "this": {}, "throw": {},
	"throws": {}, "transient": {}, "try": {}, "void": {},
	"volatile": {}, "while": {},
}

func toJavaName(name string) string {
	name = strings.ReplaceAll(name, "/", ".")

	segments := strings.Split(name, ".")
	var escapes bool
	for i, s := range segments {
		if _, ok := javaKeywords[s]; ok {
			segments[i] = s + "_"
			escapes = true
		}
	}
	if escapes {
		return strings.Join(segments, ".")
	}
	return name
}

// GenerateJava writes the code into the respective ".java" files.
func GenerateJava(basedir string, packages Packages) error {
	titleCache := make(map[string]string)
	funcs := template.FuncMap{"title": func(s string) string {
		if t, ok := titleCache[s]; ok {
			return t
		}

		r, size := utf8.DecodeRuneInString(s)
		if size == 0 {
			return s
		}
		t := string([]rune{unicode.ToUpper(r)}) + s[size:]

		titleCache[s] = t
		return t
	}}

	packageTemplate := template.New("java-package")
	template.Must(packageTemplate.Parse(javaPackage))
	codeTemplate := template.New("java-code").Funcs(funcs)
	template.Must(codeTemplate.Parse(javaCode))

	for _, p := range packages {
		p.NameNative = toJavaName(p.Name)
		p.SuperClassNative = toJavaName(p.SuperClass)

		p.InterfaceNatives = make([]string, len(p.Interfaces))
		for i, s := range p.Interfaces {
			p.InterfaceNatives[i] = toJavaName(s)
		}

		for _, t := range p.Structs {
			t.NameNative = name.CamelCase(t.Name, true)
			for _, f := range t.Fields {
				f.NameNative = name.CamelCase(f.Name, false)
				if _, ok := javaKeywords[f.NameNative]; ok {
					f.NameNative += "_"
				}
			}
		}
	}

	for _, p := range packages {
		pkgdir := filepath.Join(basedir, strings.Replace(p.NameNative, ".", string([]rune{filepath.Separator}), -1))
		if err := os.MkdirAll(pkgdir, os.ModeDir|os.ModePerm); err != nil {
			return err
		}

		if doc := p.DocText(" * "); doc != "" {
			f, err := os.Create(filepath.Join(pkgdir, "package-info.java"))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := packageTemplate.Execute(f, p); err != nil {
				return err
			}
		}

		for _, t := range p.Structs {
			for _, f := range t.Fields {
				switch f.Type {
				default:
					if f.TypeRef == nil {
						f.TypeNative = f.Type
					} else {
						f.TypeNative = f.TypeRef.NameNative
						if f.TypeRef.Pkg != p {
							f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeNative
						}
					}
				case "bool":
					f.TypeNative = "boolean"
				case "uint8":
					f.TypeNative = "byte"
				case "uint16":
					f.TypeNative = "short"
				case "uint32", "int32":
					f.TypeNative = "int"
				case "uint64", "int64":
					f.TypeNative = "long"
				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"
				case "timestamp":
					f.TypeNative = "java.time.Instant"
				case "text":
					f.TypeNative = "String"
				case "binary":
					f.TypeNative = "byte[]"
				}
			}

			f, err := os.Create(filepath.Join(pkgdir, t.NameNative+".java"))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := codeTemplate.Execute(f, t); err != nil {
				return err
			}
		}
	}
	return nil
}

//go:embed template/java-package.txt
var javaPackage string

//go:embed template/java.txt
var javaCode string
