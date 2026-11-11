//go:build !windows

package helpers

import (
	"os"
	"syscall"
)

// requestWebUIShutdown asks the Web UI daemon to terminate gracefully.
func requestWebUIShutdown(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
