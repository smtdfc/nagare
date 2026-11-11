package system

import (
	"fmt"
)

// BrightnessControl adjusts the display brightness of the primary active
// display as a percentage between 0 and 100.
//
// Backends per OS:
//   - Linux/BSD: brightnessctl (works under both X11 and Wayland)
//   - Windows: WMI WmiMonitorBrightnessMethods via PowerShell
//   - macOS: CoreGraphics CGDisplaySetBrightness via purego (no cgo)
type BrightnessControl struct{}

func NewBrightnessControl() *BrightnessControl {
	return &BrightnessControl{}
}

// Get returns the current brightness as a percentage in [0, 100].
func (b *BrightnessControl) Get() (int, error) {
	return brightnessGet()
}

// Set applies a brightness percentage in [0, 100].
func (b *BrightnessControl) Set(value int) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("brightness value must be between 0 and 100")
	}
	return brightnessSet(value)
}

// Increment raises the brightness by the given delta, clamped to 100.
func (b *BrightnessControl) Increment(delta int) (int, error) {
	current, err := b.Get()
	if err != nil {
		return 0, err
	}
	next := clampPercent(current + delta)
	return next, b.Set(next)
}

// Decrement lowers the brightness by the given delta, clamped to 0.
func (b *BrightnessControl) Decrement(delta int) (int, error) {
	current, err := b.Get()
	if err != nil {
		return 0, err
	}
	next := clampPercent(current - delta)
	return next, b.Set(next)
}

func clampPercent(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}
