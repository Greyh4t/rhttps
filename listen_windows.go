package main

import (
	"context"
	"net"
	"syscall"

	"golang.org/x/sys/windows"
)

func Listen(network, addr string, reusePort bool) (net.Listener, error) {
	if !reusePort {
		return net.Listen(network, addr)
	}

	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) (err error) {
			if err := c.Control(func(fd uintptr) {
				err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
			}); err != nil {
				return err
			}
			return
		},
	}
	return lc.Listen(context.Background(), network, addr)
}
