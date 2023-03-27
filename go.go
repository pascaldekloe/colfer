package colfer

import (
	_ "embed"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pascaldekloe/name"
	"golang.org/x/mod/modfile"
)

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
