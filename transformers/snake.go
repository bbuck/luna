// Copyright (c) 2020 Brandon Buck

package transformers

import (
	"unicode"
)

// StringToSnake converts an exported Go name to snake_case following the Golang
// format: acronyms are converted to lower-case and preceded by an underscore.
// found at: https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func StringToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
