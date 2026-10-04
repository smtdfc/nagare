//go:build !windows

package helpers

import (
	"os"
	"syscall"
)

func requestGatewayShutdown(process *os.Process) error {
	return process.Signal(syscall.SIGTERM)
}
