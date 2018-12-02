// +build !linux !386
// +build !linux !amd64
// +build !linux !arm
// +build !solaris
// +build !windows

package netutil

import "syscall"

func listenControl(network, address string, c syscall.RawConn) (err error) {
	return c.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		if err != nil {
			return
		}
		err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
	})
}
