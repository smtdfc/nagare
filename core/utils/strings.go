package utils

import "unicode/utf8"

func TruncateString(str string, maxChars int) string {
	if utf8.RuneCountInString(str) <= maxChars {
		return str
	}

	runes := []rune(str)
	return string(runes[:maxChars]) + "... (truncated)"
}
