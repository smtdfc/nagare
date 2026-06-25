package media

import (
	"fmt"

	"github.com/itchyny/volume-go"
)

func SetVolume(percent int) error {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	if percent == 0 {
		return volume.Mute()
	}

	isMuted, err := volume.GetMuted()
	if err != nil {
		return err
	}

	if isMuted {
		if err := volume.Unmute(); err != nil {
			return err
		}
	}

	if err := volume.SetVolume(percent); err != nil {
		return fmt.Errorf("failed to set volume: %w", err)
	}

	return nil
}
