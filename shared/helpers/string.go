package helpers

import "strings"

func ContainsAnyKeyword(target string, keywords []string) bool {
	targetLower := strings.ToLower(target)
	for _, kw := range keywords {
		if strings.Contains(targetLower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
