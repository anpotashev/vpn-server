package ifaceconfigurator

import (
	"fmt"
	"net"
	"os"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

type impl struct{}

func New() Configurator {
	return &impl{}
}

func (m *impl) InitIface(ip net.IP, mask net.IPMask, mtu uint16) (*water.Interface, func() error, error) {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, nil, err
	}
	name := iface.Name()
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, nil, err
	}
	addr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   ip,
			Mask: mask,
		},
		Peer: &net.IPNet{
			IP:   net.IPv4(192, 168, 100, 2), // IP клиента
			Mask: net.CIDRMask(32, 32),
		},
	}
	err = netlink.AddrAdd(link, addr)
	if err != nil {
		return nil, nil, err
	}
	err = netlink.LinkSetMTU(link, int(mtu))
	if err != nil {
		return nil, nil, err
	}
	err = netlink.LinkSetUp(link)
	if err != nil {
		return nil, nil, err
	}
	os.WriteFile(fmt.Sprintf("/proc/sys/net/ipv4/conf/%s/forwarding", iface.Name()), []byte("1"), 0644)
	return iface, nil, nil
}
