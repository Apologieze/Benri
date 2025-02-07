//go:build windows
// +build windows

package curdInteg

import (
	"github.com/Microsoft/go-winio"
	"net"
)

func connectToPipe(ipcSocketPath string) (net.Conn, error) {
	conn, err := winio.DialPipe(ipcSocketPath, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
