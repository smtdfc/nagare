//go:build windows

package helpers

import "os"

// requestWebUIShutdown asks the Web UI daemon to terminate. Windows processes
// only support the hard kill signal through os.Process.
func requestWebUIShutdown(process *os.Process) error {
	return process.Signal(os.Interrupt)
}
