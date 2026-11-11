package cmd

import "testing"

func TestWebStopCommandRegistered(t *testing.T) {
	stop, _, err := webCmd.Find([]string{"stop"})
	if err != nil {
		t.Fatalf("unexpected error resolving 'nagare web stop': %v", err)
	}
	if stop == nil || stop.Name() != "stop" {
		t.Fatalf("expected 'nagare web stop' to resolve to the stop command, got %+v", stop)
	}
}
