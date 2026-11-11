//go:build darwin

package system

import (
	"fmt"

	"github.com/ebitengine/purego"
)

// cgLib holds the CoreGraphics entry points used for brightness control.
type cgLib struct {
	displayBrightness    func(display uint32) float64
	setDisplayBrightness func(display uint32, level float32) int32
	getActiveDisplayList func(maxDisplays uint32, activeDisplays *uint32, displayCount *uint32) int32
}

var coreGraphicsLib *cgLib

// loadCoreGraphics resolves the CoreGraphics symbols lazily through purego so
// the code stays cgo-free and only pays the cost on a real macOS host.
func loadCoreGraphics() (*cgLib, error) {
	if coreGraphicsLib != nil {
		return coreGraphicsLib, nil
	}

	lib, err := purego.Dlopen("/System/Library/Frameworks/CoreGraphics.framework/CoreGraphics", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("failed to load CoreGraphics: %w", err)
	}

	var control cgLib
	purego.RegisterLibFunc(&control.displayBrightness, lib, "CGDisplayBrightness")
	purego.RegisterLibFunc(&control.setDisplayBrightness, lib, "CGDisplaySetBrightness")
	purego.RegisterLibFunc(&control.getActiveDisplayList, lib, "CGGetActiveDisplayList")

	coreGraphicsLib = &control
	return coreGraphicsLib, nil
}

// activeMacDisplay returns the first active display identifier.
func activeMacDisplay(control *cgLib) (uint32, error) {
	var displays [8]uint32
	var count uint32
	if ret := control.getActiveDisplayList(uint32(len(displays)), &displays[0], &count); ret != 0 || count == 0 {
		return 0, fmt.Errorf("failed to resolve the active display (ret=%d count=%d)", ret, count)
	}
	return displays[0], nil
}

func brightnessGet() (int, error) {
	control, err := loadCoreGraphics()
	if err != nil {
		return 0, err
	}
	display, err := activeMacDisplay(control)
	if err != nil {
		return 0, err
	}
	return int(control.displayBrightness(display) * 100), nil
}

func brightnessSet(value int) error {
	control, err := loadCoreGraphics()
	if err != nil {
		return err
	}
	display, err := activeMacDisplay(control)
	if err != nil {
		return err
	}
	if ret := control.setDisplayBrightness(display, float32(value)/100); ret != 0 {
		return fmt.Errorf("CGDisplaySetBrightness failed with %d", ret)
	}
	return nil
}
