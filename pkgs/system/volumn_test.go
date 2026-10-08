package system

import (
	"testing"
)

func TestVolumeControlInstance(t *testing.T) {
	vc := NewVolumeControl()
	if vc == nil {
		t.Fatalf("expected non-nil VolumeControl instance")
	}

	// In CI/headless environments, audio server (e.g. pulse/alsa) might not be running.
	// We test method dispatch without panicking.
	vol, err := vc.Get()
	if err != nil {
		t.Logf("volume control not available in this environment: %v", err)
	} else {
		t.Logf("current volume: %d", vol)
	}
}
