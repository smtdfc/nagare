package hooks

import "github.com/google/uuid"

var onStartHooks = make(map[string]HookFn)
var onStopHooks = make(map[string]HookFn)

func OnStart(fn HookFn) string {
	hookID := uuid.New().String()
	onStartHooks[hookID] = fn
	return hookID
}

func OnStop(fn HookFn) string {
	hookID := uuid.New().String()
	onStopHooks[hookID] = fn
	return hookID
}

func GetOnStartHooks() map[string]HookFn {
	return onStartHooks
}

func GetOnStopHooks() map[string]HookFn {
	return onStopHooks
}
