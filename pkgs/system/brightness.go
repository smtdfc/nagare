package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/ebitengine/purego"
)

// BrightnessControl adjusts the display brightness of the primary active
// display as a percentage between 0 and 100.
//
// Backends per OS:
//   - Linux: brightnessctl (works under both X11 and Wayland, rootless via
//     uaccess rules or sudoer setups)
//   - Windows: WMI WmiMonitorBrightnessMethods via PowerShell
//   - macOS: CoreGraphics CGDisplaySetBrightness via purego (no cgo)
type BrightnessControl struct{}

func NewBrightnessControl() *BrightnessControl {
	return &BrightnessControl{}
}

// Get returns the current brightness as a percentage in [0, 100].
func (b *BrightnessControl) Get() (int, error) {
	switch runtime.GOOS {
	case "linux":
		return brightnessGetLinux()
	case "windows":
		return brightnessGetWindows()
	case "darwin":
		return brightnessGetDarwin()
	default:
		return 0, fmt.Errorf("brightness control is not supported on %s", runtime.GOOS)
	}
}

// Set applies a brightness percentage in [0, 100].
func (b *BrightnessControl) Set(value int) error {
	if value < 0 || value > 100 {
		return fmt.Errorf("brightness value must be between 0 and 100")
	}

	switch runtime.GOOS {
	case "linux":
		return brightnessSetLinux(value)
	case "windows":
		return brightnessSetWindows(value)
	case "darwin":
		return brightnessSetDarwin(value)
	default:
		return fmt.Errorf("brightness control is not supported on %s", runtime.GOOS)
	}
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

// --- Linux ---

func brightnessGetLinux() (int, error) {
	out, err := exec.Command("brightnessctl", "info", "--machine-readable").Output()
	if err != nil {
		return 0, fmt.Errorf("brightnessctl is required for brightness control on Linux (install with your package manager): %w", err)
	}

	// Machine-readable format: <name>,<device>,<current>,<max>,<percent>,...
	fields := strings.Split(strings.TrimSpace(string(out)), ",")
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if strings.HasSuffix(trimmed, "%") {
			percent, err := strconv.Atoi(strings.TrimSuffix(trimmed, "%"))
			if err != nil {
				return 0, fmt.Errorf("failed to parse brightnessctl output %q: %w", trimmed, err)
			}
			return percent, nil
		}
	}

	return 0, fmt.Errorf("failed to parse brightnessctl output: %q", string(out))
}

func brightnessSetLinux(value int) error {
	cmd := exec.Command("brightnessctl", "set", "--quiet", fmt.Sprintf("%d%%", value))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set brightness via brightnessctl: %w", err)
	}
	return nil
}

// --- Windows ---

func brightnessGetWindows() (int, error) {
	out, err := runPowerShell("(Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightness).CurrentBrightness")
	if err != nil {
		return 0, fmt.Errorf("failed to read brightness (desktop WMI monitors only): %w", err)
	}
	value, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("failed to parse brightness output %q: %w", out, err)
	}
	return value, nil
}

func brightnessSetWindows(value int) error {
	_, err := runPowerShell(fmt.Sprintf(
		"Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightnessMethods | Invoke-CimMethod -MethodName WmiSetBrightness -Arguments @{Timeout=1; Brightness=%d}",
		value,
	))
	if err != nil {
		return fmt.Errorf("failed to set brightness (desktop WMI monitors only): %w", err)
	}
	return nil
}

// --- macOS ---

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

func brightnessGetDarwin() (int, error) {
	control, err := loadCoreGraphics()
	if err != nil {
		return 0, err
	}
	display, err := activeMacDisplay(control)
	if err != nil {
		return 0, err
	}
	level := control.displayBrightness(display)
	return int(level * 100), nil
}

func brightnessSetDarwin(value int) error {
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

// runPowerShell executes a single PowerShell expression and captures stdout.
func runPowerShell(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(out), err
}
