package tui

import "net"

type networkInterfacesMsg []net.Interface

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }
