// Copyright (c) 2020 Brandon Buck

package transformers

import "unicode"

// StringToCamel converts an exported Go name to camelCase.
func StringToCamel(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	if length > 1 {
		switch {
		case unicode.IsUpper(runes[0]) && unicode.IsLower(runes[1]):
			out = append(out, unicode.ToLower(runes[0]))
			out = append(out, runes[1:]...)
		case unicode.IsUpper(runes[0]) && unicode.IsUpper(runes[1]):
			i := 0
			for ; i < length && unicode.IsUpper(runes[i]); i++ {
				if i+1 < length && unicode.IsUpper(runes[i+1]) {
					out = append(out, unicode.ToLower(runes[i]))
				} else {
					out = append(out, runes[i])
				}
			}
			if i < length {
				out = append(out, runes[i:]...)
			}
		default:
			out = append(out, runes[0:]...)
		}
	} else if length == 1 {
		out = append(out, unicode.ToLower(runes[0]))
	}

	return string(out)
}
