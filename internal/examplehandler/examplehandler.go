package examplehandler

import (
	"log/slog"
	"net"
	"net/http"

	"github.com/anpotashev/vpn-server/pkg/vpnproxy"
	"github.com/gorilla/websocket"
)

type exampleHandler struct {
	v vpnproxy.VPNProxy
}

func NewExampleHandler() (*exampleHandler, error) {
	proxy, err := vpnproxy.NewVPNProxy(vpnproxy.VPNProxyConfig{
		IP:         net.IPv4(192, 168, 100, 1),
		IPMask:     net.IPv4Mask(255, 255, 255, 255),
		MTU:        1500,
		IPTemplate: []byte{192, 168, 100},
	})
	if err != nil {
		return nil, err
	}
	return &exampleHandler{proxy}, nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (e *exampleHandler) ExampleHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	vpn := e.v.AttachVPN(c)
	defer vpn.Close()
	for {
		msgType, data, err := vpn.Read()
		if err != nil {
			slog.Error("Error reading message from the websocket", "err", err)
			break
		}
		if msgType == websocket.TextMessage && string(data) == "ping" {
			err := vpn.Write(websocket.TextMessage, []byte("pong"))
			if err != nil {
				slog.Error("Error writing message to the websocket", "err", err)
				return
			}
		}
	}
}
