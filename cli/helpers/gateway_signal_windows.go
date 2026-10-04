//go:build windows

package helpers

import "os"

func requestGatewayShutdown(process *os.Process) error {
	return process.Signal(os.Interrupt)
}
