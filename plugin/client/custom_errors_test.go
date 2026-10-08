package client

import (
	"testing"
)

func TestCustomErrors(t *testing.T) {
	errs := []error{
		ErrConnectionFailed,
		ErrConnectionNotFound,
		ErrConnectionNotReady,
		ErrCreatePayloadFailed,
		ErrHandshakeFailed,
		ErrHandshakeTimeout,
		ErrHandshakeCancelled,
		ErrPrepareChatSessionFailed,
		ErrIncorrectToolArgs,
		ErrMarshalToolResultFailed,
		ErrTimeout,
		ErrActionCancelled,
	}

	for _, err := range errs {
		if err == nil || err.Error() == "" {
			t.Fatalf("expected non-empty error message, got %v", err)
		}
	}
}
