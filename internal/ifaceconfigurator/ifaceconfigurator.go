package ifaceconfigurator

import (
	"errors"
	"net"

	"github.com/songgao/water"
)

var ErrNotImplemented = errors.New("not implemented")

type Configurator interface {
	InitIface(ip net.IP, mask net.IPMask, mtu uint16) (*water.Interface, func() error, error)
}
