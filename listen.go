//go:build linux || darwin

package main

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func Listen(network, addr string, reusePort bool) (net.Listener, error) {
	if !reusePort {
		return net.Listen(network, addr)
	}

	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) (err error) {
			if err := c.Control(func(fd uintptr) {
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				if err != nil {
					return
				}
				err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return
		},
	}
	return lc.Listen(context.Background(), network, addr)
}
