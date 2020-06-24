package colfer

import "testing"

func GoldenTagPackages() Packages {
	p := &Package{Name: "gen"}
	t := &Struct{Name: "o", Pkg: p}
	f := &Field{Name: "i", Struct: t}

	p.Structs = append(p.Structs, t)
	t.Fields = append(t.Fields, f)

	return Packages{p}
}

var GoldenTagFileErrors = []struct{ File, Err string }{
	{
		"testdata/struct-package-miss.tags",
		`testdata/struct-package-miss.tags:2: "wrong.o" package not in schema`,
	}, {
		"testdata/struct-miss.tags",
		`testdata/struct-miss.tags:2: "gen.wrong" struct not in schema`,
	}, {
		"testdata/field-package-miss.tags",
		`testdata/field-package-miss.tags:2: "wrong.o.i" package not in schema`,
	}, {
		"testdata/field-struct-miss.tags",
		`testdata/field-struct-miss.tags:2: "gen.wrong.i" struct not in schema`,
	}, {
		"testdata/field-miss.tags",
		`testdata/field-miss.tags:2: "gen.o.wrong" field not in schema`,
	}, {
		"testdata/package-case.tags",
		`testdata/package-case.tags:2: "Gen.o.i" case mismatch with gen?`,
	}, {
		"testdata/struct-case.tags",
		`testdata/struct-case.tags:2: "gen.O.i" case mismatch with gen.o?`,
	}, {
		"testdata/field-case.tags",
		`testdata/field-case.tags:2: "gen.o.I" case mismatch with gen.o.i?`,
	}, {
		"testdata/corrupt.tags",
		`testdata/corrupt.tags:6: invalid qualifier "broken"; use <package>.<type> or <package>.<type>.<field>`,
	},
}

func TestTagFileErrors(t *testing.T) {
	for _, gold := range GoldenTagFileErrors {
		err := GoldenTagPackages().ApplyTagFile(gold.File)
		if err == nil || err.Error() != gold.Err {
			t.Errorf("%s: got %v, want %v", gold.File, err, gold.Err)
		}
	}
}
