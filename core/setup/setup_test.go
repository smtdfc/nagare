package setup

import (
	"log/slog"
	"testing"

	"github.com/smtdfc/nagare/core/logger"
)

func TestCoreSetup_Struct(t *testing.T) {
	// Verify CoreSetup constructor initialization
	bl := &logger.BaseLogger{Logger: *slog.Default()}
	cs := NewCoreSetup(nil, nil, nil, nil, nil, nil, nil, bl)

	if cs == nil {
		t.Fatalf("expected non-nil CoreSetup")
	}
	if cs.logger == nil {
		t.Errorf("expected logger to be initialized")
	}
}
