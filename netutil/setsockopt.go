package netutil

import "syscall"

func ListenControl(network, address string, c syscall.RawConn) error {
	return listenControl(network, address, c)
}
