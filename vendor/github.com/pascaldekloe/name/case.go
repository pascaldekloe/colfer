// Package name implements naming conventions like camel and snake case.
package name

import "unicode"

// CamelCase returns the camel case of word sequence s.
// The input can be any case or even just a bunch of words.
// Upper case abbreviations are preserved. When upper is true
// then the first rune is mapped to its upper case form.
func CamelCase(s string, upper bool) string {
	out := make([]rune, 0, len(s)+5)
	for i, r := range s {
		switch {
		case i == 0:
			if !upper {
				r = unicode.ToLower(r)
			}

			fallthrough
		case unicode.IsLetter(r):
			if upper {
				r = unicode.ToUpper(r)
			}

			fallthrough
		case unicode.IsNumber(r):
			upper = false
			out = append(out, r)

		default:
			upper = true

		}
	}
	return string(out)
}

// SnakeCase returns the snake case of word sequence s.
// The input can be any case or even just a bunch of words.
// Upper case abbreviations are preserved. Use strings.ToLower
// and strings.ToUpper to enforce a letter case.
func SnakeCase(s string) string {
	return Delimit(s, '_')
}

// Delimit returns word sequence s delimited with sep.
// The input can be any case or even just a bunch of words.
// Upper case abbreviations are preserved. Use strings.ToLower
// and strings.ToUpper to enforce a letter case.
func Delimit(s string, sep rune) string {
	out := make([]rune, 0, len(s)+5)

	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			if last := len(out) - 1; last >= 0 && unicode.IsLower(out[last]) {
				out = append(out, sep)
			}

		case unicode.IsLetter(r):
			if i := len(out) - 1; i >= 0 {
				if last := out[i]; unicode.IsUpper(last) {
					out = out[:i]
					if i > 0 && out[i-1] != sep {
						out = append(out, sep)
					}
					out = append(out, unicode.ToLower(last))
				}
			}

		case !unicode.IsNumber(r):
			if i := len(out); i != 0 && out[i-1] != sep {
				out = append(out, sep)
			}
			continue

		}
		out = append(out, r)
	}

	// trim tailing separator
	if i := len(out) - 1; i >= 0 && out[i] == sep {
		out = out[:i]
	}

	return string(out)
}
