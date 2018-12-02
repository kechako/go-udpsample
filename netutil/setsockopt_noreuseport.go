// +build linux,386 linux,amd64 linux,arm solaris

package netutil

import (
	"syscall"
)

func listenControl(network, address string, c syscall.RawConn) (err error) {
	return c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
}
