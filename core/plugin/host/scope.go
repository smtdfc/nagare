package host

import (
	"slices"

	"github.com/smtdfc/nagare/core/plugin"
)

func RequireScopes(raws []string, scopes ...plugin.Feature) bool {
	for _, scope := range scopes {
		if !slices.Contains(raws, scope.ToString()) {
			return false
		}
	}

	return true
}
