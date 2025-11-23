package vpnproxy

import (
	"log/slog"
	"net"

	"github.com/anpotashev/vpn-server/internal/ifaceconfigurator"
	"github.com/anpotashev/vpn-server/internal/ipallocator"
	"github.com/gorilla/websocket"
	"github.com/songgao/water"
)

type AppConn interface {
	Read() (messageType int, data []byte, err error) // только не-VPN данные
	Write(messageType int, data []byte) error        // отправка не-VPN данных клиенту
	Close() error
}

type VPNProxy interface {
	// AttachVPN перехватывает сообщения полученные от WS и обрабатывает только относящиеся к VPN,
	// остальные отдает на обработку основному приложению
	// основное приложение не должно вызывать методы чтения у переданного conn, только этот.
	AttachVPN(conn *websocket.Conn) AppConn
}

type VPNProxyConfig struct {
	IPTemplate []byte
	//iface      *water.Interface
	MTU    uint16
	IP     net.IP
	IPMask net.IPMask
}

type vpnProxy struct {
	ipAllocator ipallocator.IPAllocator
	iface       *water.Interface
	clients     map[*client]struct{}
	mtu         uint16
}

func NewVPNProxy(config VPNProxyConfig) (VPNProxy, error) {
	iface, _, err := ifaceconfigurator.New().InitIface(config.IP, config.IPMask, config.MTU)
	if err != nil {
		return nil, err
	}
	proxy := &vpnProxy{
		//ws:          conn,
		ipAllocator: ipallocator.NewIP4Allocator(config.IPTemplate),
		iface:       iface,
		clients:     make(map[*client]struct{}),
		mtu:         config.MTU,
	}
	go proxy.startListeningInterface()
	return proxy, nil
}

func (v *vpnProxy) startListeningInterface() {
	buf := make([]byte, 65535)
	for {
		n, _ := v.iface.Read(buf)
		slog.Debug("Received packet from the interface", headerDebugData(buf[:n])...)
		pkg := buf[:n]
		dst := parseIPv4Dst(pkg)
		if dst == nil {
			continue
		}
		for c := range v.clients {
			if c.ip.String() == dst.String() {
				_ = c.conn.WriteMessage(websocket.BinaryMessage, pkg)
				break
			}
		}
	}
}

func (v *vpnProxy) AttachVPN(conn *websocket.Conn) AppConn {
	c := &client{
		conn:        conn,
		ipAllocator: v.ipAllocator,
		mtu:         v.mtu,
		writeToInterface: func(payload []byte) error {
			slog.Debug("Sending the packet to the interface:", headerDebugData(payload)...)
			_, err := v.iface.Write(payload)
			return err
		},
	}
	v.clients[c] = struct{}{}
	c.beforeClose = func() { delete(v.clients, c) }
	return c
}
