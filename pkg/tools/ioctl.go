package tools

import "syscall"

func Ioctl(a1, a2, a3 uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, a1, a2, a3)
	if errno != 0 {
		return errno
	}
	return nil
}
