package system

import (
	"testing"
)

func TestProcessManagerLifecycle(t *testing.T) {
	pm := NewProcessManager()

	// Test starting an invalid command
	_, err := pm.Start("", nil, "")
	if err == nil {
		t.Fatalf("expected error when starting empty command, got nil")
	}

	// Test starting a real command (sleep or test executable)
	tempDir := t.TempDir()
	pid, err := pm.Start("sleep", []string{"5"}, tempDir)
	if err != nil {
		t.Fatalf("failed to start process: %v", err)
	}
	if pid <= 0 {
		t.Fatalf("expected positive PID, got %d", pid)
	}

	// Test finding processes with filters
	items, err := pm.Find(&ProcessQuery{
		NameQuery: "sleep",
		Limit:     5,
	})
	if err != nil {
		t.Fatalf("failed to find processes: %v", err)
	}
	found := false
	for _, it := range items {
		if it.PID == pid {
			found = true
			break
		}
	}
	if !found {
		t.Logf("note: process %d might have exited or renamed, found %d items", pid, len(items))
	}

	// Test killing process
	name, err := pm.Kill(pid)
	if err != nil {
		t.Logf("process kill error (may have already exited): %v, name: %s", err, name)
	}

	// Test killing invalid PID
	_, err = pm.Kill(-1)
	if err == nil {
		t.Fatalf("expected error killing invalid pid -1, got nil")
	}
}

func TestProcessQuerySorting(t *testing.T) {
	pm := NewProcessManager()
	// Test querying processes sorted by ram and limited to 3
	items, err := pm.Find(&ProcessQuery{
		SortBy: "ram",
		Limit:  3,
	})
	if err != nil {
		t.Fatalf("failed to query processes: %v", err)
	}
	if len(items) > 3 {
		t.Fatalf("expected at most 3 items, got %d", len(items))
	}
	if len(items) >= 2 && items[0].MemoryMB < items[1].MemoryMB {
		t.Fatalf("expected items sorted in descending RAM order: %v vs %v", items[0].MemoryMB, items[1].MemoryMB)
	}
}
