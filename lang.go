package colfer

import (
	"bytes"
	_ "embed"
	"io/ioutil"
	"os"
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
	// set native names and types
	for _, p := range packages {
		for _, t := range p.Structs {
			t.NameNative = strings.ToLower(name.SnakeCase(p.Name + "_" + t.Name))

			for _, f := range t.Fields {
				if f.Type == "bool" {
					fqn := name.SnakeCase(p.Name + "_" + t.Name + "_" + f.Name)
					f.NameNative = strings.ToUpper(fqn) + "_FLAG"
				} else {
					f.NameNative = strings.ToLower(name.SnakeCase(f.Name))
					if _, ok := cKeywords[f.NameNative]; ok {
						f.NameNative += "_"
					}
				}

				switch f.Type {
				case "opaque8", "uint8":
					f.TypeNative = "uint8_t"
				case "opaque16", "uint16":
					f.TypeNative = "uint16_t"
				case "opaque32", "uint32":
					f.TypeNative = "uint32_t"
				case "opaque64", "uint64":
					f.TypeNative = "uint64_t"

				case "int8":
					f.TypeNative = "int8_t"
				case "int16":
					f.TypeNative = "int16_t"
				case "int32":
					f.TypeNative = "int32_t"
				case "int64":
					f.TypeNative = "int64_t"

				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"

				case "timestamp":
					f.TypeNative = "struct timespec"
				case "text":
					f.TypeNative = "char*"
				case "bool":
					f.TypeNative = "int"
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

	funcs := template.FuncMap{"upperSnake": upperSnake}
	headerTemplate, err := template.New("c-header").Funcs(funcs).Parse(hTemplate)
	if err != nil {
		return err
	}
	if err := headerTemplate.Execute(f, packages); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	f, err = os.Create(filepath.Join(basedir, "Colfer.c"))
	if err != nil {
		return err
	}
	t := template.New("c")
	template.Must(t.Parse(cTemplate))
	template.Must(t.New("marshal-integer").Parse(cMarshalIntegerTemplate))
	template.Must(t.New("marshal16").Parse(cMarshal16Template))
	template.Must(t.New("marshal32").Parse(cMarshal32Template))
	template.Must(t.New("marshal64").Parse(cMarshal64Template))
	template.Must(t.New("unmarshal-integer").Parse(cUnmarshalIntegerTemplate))
	template.Must(t.New("unmarshal16").Parse(cUnmarshal16Template))
	template.Must(t.New("unmarshal32").Parse(cUnmarshal32Template))
	template.Must(t.New("unmarshal64").Parse(cUnmarshal64Template))
	template.Must(t.New("malloc-error").Parse(cMallocErrorTemplate))
	if err := t.Execute(f, packages); err != nil {
		return err
	}
	return f.Close()
}

//go:embed template/c.txt
var cTemplate string

//go:embed template/h.txt
var hTemplate string

//go:embed template/c-marshal-integer.txt
var cMarshalIntegerTemplate string

//go:embed template/c-marshal16.txt
var cMarshal16Template string

//go:embed template/c-marshal32.txt
var cMarshal32Template string

//go:embed template/c-marshal64.txt
var cMarshal64Template string

//go:embed template/c-unmarshal-integer.txt
var cUnmarshalIntegerTemplate string

//go:embed template/c-unmarshal16.txt
var cUnmarshal16Template string

//go:embed template/c-unmarshal32.txt
var cUnmarshal32Template string

//go:embed template/c-unmarshal64.txt
var cUnmarshal64Template string

//go:embed template/c-malloc-error.txt
var cMallocErrorTemplate string

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
	funcs := template.FuncMap{
		"title":      title,
		"upperSnake": upperSnake,
	}

	packageTemplate := template.New("java-package")
	template.Must(packageTemplate.Parse(javaPackageTemplate))

	classTemplate := template.New("java-class").Funcs(funcs)
	template.Must(classTemplate.Parse(javaTemplate))
	template.Must(classTemplate.New("marshal-integer").Parse(javaMarshalIntegerTemplate))
	template.Must(classTemplate.New("marshal16").Parse(javaMarshal16Template))
	template.Must(classTemplate.New("marshal32").Parse(javaMarshal32Template))
	template.Must(classTemplate.New("marshal64").Parse(javaMarshal64Template))
	template.Must(classTemplate.New("unmarshal-integer").Parse(javaUnmarshalIntegerTemplate))
	template.Must(classTemplate.New("unmarshal16").Parse(javaUnmarshal16Template))
	template.Must(classTemplate.New("unmarshal32").Parse(javaUnmarshal32Template))
	template.Must(classTemplate.New("unmarshal64").Parse(javaUnmarshal64Template))

	// set Java names
	for _, p := range packages {
		p.NameNative = toJavaName(p.Name)

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

	// set Type names (needs Java names for TypeRef)
	for _, p := range packages {
		for _, t := range p.Structs {
			for _, f := range t.Fields {
				if f.TypeRef != nil {
					// fully qualified names prevent java.lang collisions
					f.TypeNative = f.TypeRef.Pkg.NameNative + "." + f.TypeRef.NameNative

					continue
				}

				switch f.Type {
				case "opaque8", "uint8", "int8":
					f.TypeNative = "byte"
				case "opaque16", "uint16", "int16":
					f.TypeNative = "short"
				case "opaque32", "uint32", "int32":
					f.TypeNative = "int"
				case "opaque64", "uint64", "int64":
					f.TypeNative = "long"

				case "float32":
					f.TypeNative = "float"
				case "float64":
					f.TypeNative = "double"

				case "timestamp":
					f.TypeNative = "Instant"
				case "text":
					f.TypeNative = "String"
				case "bool":
					f.TypeNative = "int"
				}
			}
		}
	}

	// construct files
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
			f, err := os.Create(filepath.Join(pkgdir, t.NameNative+".java"))
			if err != nil {
				return err
			}
			defer f.Close()

			if err := classTemplate.Execute(f, t); err != nil {
				return err
			}
		}
	}

	return nil
}

//go:embed template/java-package.txt
var javaPackageTemplate string

//go:embed template/java.txt
var javaTemplate string

//go:embed template/java-marshal-integer.txt
var javaMarshalIntegerTemplate string

//go:embed template/java-marshal16.txt
var javaMarshal16Template string

//go:embed template/java-marshal32.txt
var javaMarshal32Template string

//go:embed template/java-marshal64.txt
var javaMarshal64Template string

//go:embed template/java-unmarshal-integer.txt
var javaUnmarshalIntegerTemplate string

//go:embed template/java-unmarshal16.txt
var javaUnmarshal16Template string

//go:embed template/java-unmarshal32.txt
var javaUnmarshal32Template string

//go:embed template/java-unmarshal64.txt
var javaUnmarshal64Template string

var titleCache map[string]string

func title(s string) string {
	if t, ok := titleCache[s]; ok {
		return t
	}
	if titleCache == nil {
		titleCache = make(map[string]string)
	}

	r, size := utf8.DecodeRuneInString(s)
	t := string([]rune{unicode.ToUpper(r)}) + s[size:]
	titleCache[s] = t
	return t
}

var upperSnakeCache map[string]string

func upperSnake(s string) string {
	if t, ok := upperSnakeCache[s]; ok {
		return t
	}
	if upperSnakeCache == nil {
		upperSnakeCache = make(map[string]string)
	}

	t := strings.ToUpper(name.SnakeCase(s))
	upperSnakeCache[s] = t
	return t
}
