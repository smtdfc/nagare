package logger

import (
	"log/slog"
	"testing"
)

func TestBaseLogger_Methods(t *testing.T) {
	// Verify BaseLogger With and Clone methods
	bl := &BaseLogger{
		Logger: *slog.Default(),
	}

	withLogger := bl.With("key", "val")
	if withLogger == nil {
		t.Fatalf("expected non-nil withLogger")
	}

	cloned := bl.Clone()
	if cloned == nil {
		t.Fatalf("expected non-nil cloned logger")
	}
}

func TestNewBaseLogger(t *testing.T) {
	// Verify NewBaseLogger constructor
	l, err := NewBaseLogger()
	if err != nil {
		t.Fatalf("unexpected error creating base logger: %v", err)
	}
	if l == nil {
		t.Fatalf("expected non-nil base logger")
	}
}
