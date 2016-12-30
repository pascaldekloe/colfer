// Package name implements naming conventions like camel case and snake case.
package name

import "unicode"

// CamelCase returns the camel case of word sequence s.
// The input can be any case or just a bunch of words.
// Upper case abbreviations are preserved. Use strings.ToLower,
// strings.ToUpper and strings.Title to enforce a letter case.
func CamelCase(s string) string {
	out := make([]rune, 0, len(s)+5)
	var upper bool
	for _, r := range s {
		switch {
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
			continue

		}
	}
	return string(out)
}

// SnakeCase returns the snake case of word sequence s.
// The input can be any case or just a bunch of words.
// Upper case abbreviations are preserved. Use strings.ToLower and
// strings.ToUpper to enforce a letter case.
func SnakeCase(s string) string {
	return Delimit(s, '_')
}

// Delimit returns word sequence s delimited with sep.
// The input can be any case or just a bunch of words.
// Upper case abbreviations are preserved. Use strings.ToLower and
// strings.ToUpper to enforce a letter case.
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

	if len(out) == 0 {
		return ""
	}

	// trim tailing separator
	if i := len(out) - 1; out[i] == sep {
		out = out[:i]
	}

	return string(out)
}
