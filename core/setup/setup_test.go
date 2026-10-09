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

	if cs.agentPool == nil {
		t.Errorf("expected agent pool to be initialized")
	}

	if cs.chatWorker == nil {
		t.Errorf("expected chat worker to be initialized")
	}

	if cs.pluginHost == nil {
		t.Errorf("expected plugin host to be initialized")
	}

	if cs.pluginMgr == nil {
		t.Errorf("expected plugin manager to be initialized")
	}

	if cs.taskMgr == nil {
		t.Errorf("expected task manager to be initialized")
	}

	if cs.taskCronJobWorker == nil {
		t.Errorf("expected task cron job worker to be initialized")
	}

	if cs.vectorMemory == nil {
		t.Errorf("expected vector memory to be initialized")
	}
}
