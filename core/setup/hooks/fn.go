package hooks

import (
	"context"
)

type HookFn func(context.Context) error
