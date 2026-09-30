package helpers

import (
	"strings"
)

func ContainsAnyWord(target string, keywords []string) bool {
	targetLower := strings.ToLower(target)

	targetWords := strings.Fields(targetLower)

	for _, kw := range keywords {
		kwLower := strings.ToLower(strings.TrimSpace(kw))
		if kwLower == "" {
			continue
		}

		kwParts := strings.Fields(kwLower)
		for _, part := range kwParts {
			for _, tWord := range targetWords {
				if tWord == part {
					return true
				}
			}
		}
	}
	return false
}
