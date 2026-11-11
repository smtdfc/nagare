//go:build !windows

package helpers

import "syscall"

// detachSysProcAttr makes a child process survive the CLI process: on Unix
// the child joins a new session, detaching it from the terminal signals and
// the controlling TTY.
func detachSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
