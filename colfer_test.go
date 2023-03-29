package colfer

import "testing"

func GoldenTagPackages() Packages {
	p := &Package{Name: "gen"}
	t := &Struct{Name: "o", Pkg: p}
	f := &Field{Name: "p", Struct: t}

	p.Structs = append(p.Structs, t)
	t.Fields = append(t.Fields, f)

	return Packages{p}
}

var GoldenTagFileErrors = []struct{ File, Err string }{
	{
		"testdata/struct-package-miss.tags",
		`map testdata/struct-package-miss.tags:2: package "wrong" not in schema`,
	}, {
		"testdata/struct-miss.tags",
		`map testdata/struct-miss.tags:2: type "gen.wrong" not in schema`,
	}, {
		"testdata/field-package-miss.tags",
		`map testdata/field-package-miss.tags:2: package "wrong" not in schema`,
	}, {
		"testdata/field-struct-miss.tags",
		`map testdata/field-struct-miss.tags:2: type "gen.wrong" not in schema`,
	}, {
		"testdata/field-miss.tags",
		`map testdata/field-miss.tags:2: field "gen.o.wrong" not in schema`,
	}, {
		"testdata/package-case.tags",
		`map testdata/package-case.tags:2: package not found; case mismatch with gen?`,
	}, {
		"testdata/struct-case.tags",
		`map testdata/struct-case.tags:2: type not found; case mismatch with gen.o?`,
	}, {
		"testdata/field-case.tags",
		`map testdata/field-case.tags:2: field not found; case mismatch with gen.o.p?`,
	}, {
		"testdata/corrupt.tags",
		`parse testdata/corrupt.tags:6: invalid qualifier "broken"; use <package>'.'<type>('.'<field>)`,
	},
}

func TestTagFileErrors(t *testing.T) {
	options := TagOptions{StructAllow: TagMulti, FieldAllow: TagMulti}
	for _, gold := range GoldenTagFileErrors {
		err := GoldenTagPackages().ApplyTagFile(gold.File, options)
		if err == nil || err.Error() != gold.Err {
			t.Errorf("%s: got %v, want %v", gold.File, err, gold.Err)
		}
	}
}
