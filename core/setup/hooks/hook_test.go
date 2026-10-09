package hooks

import (
	"context"
	"testing"
)

func TestOnStart(t *testing.T) {
	hookFn := func(context.Context) error {
		return nil
	}
	hookID := OnStart(hookFn)
	_, ok := onStartHooks[hookID]
	if !ok {
		t.Errorf("OnStart hook not registered")
	}
}

func TestOnStop(t *testing.T) {
	hookFn := func(context.Context) error {
		return nil
	}
	hookID := OnStop(hookFn)
	_, ok := onStopHooks[hookID]
	if !ok {
		t.Errorf("OnStop hook not registered")
	}
}
