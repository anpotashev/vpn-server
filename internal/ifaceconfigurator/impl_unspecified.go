//go:build !linux
// +build !linux

package ifaceconfigurator

import (
	"net"

	"github.com/songgao/water"
)

type impl struct{}

func New() Configurator {
	return &impl{}
}

func (m *impl) InitIface(ip net.IP, mask net.IPMask, mtu uint16) (*water.Interface, func() error, error) {
	return nil, nil, ErrNotImplemented
}
