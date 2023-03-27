package colfer

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/pascaldekloe/name"
)

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
