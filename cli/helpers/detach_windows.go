//go:build windows

package helpers

import "syscall"

// detachedProcessCreationFlag matches the Win32 DETACHED_PROCESS creation
// flag; the syscall package does not export it by name.
const detachedProcessCreationFlag = 0x00000008

// detachSysProcAttr makes a child process survive the CLI process on
// Windows: it starts in its own process group without inheriting the
// console.
func detachSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | detachedProcessCreationFlag,
	}
}
