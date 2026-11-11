//go:build windows

package helpers

import "syscall"

const detachedProcessCreationFlag = 0x00000008

func detachSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | detachedProcessCreationFlag,
	}
}
