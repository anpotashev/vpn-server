package vpnproxy

import "github.com/gorilla/websocket"

type noProxy struct {
	conn *websocket.Conn
}

func (p *noProxy) Read() (messageType int, data []byte, err error) {
	return p.conn.ReadMessage()
}

func (p *noProxy) Write(messageType int, data []byte) error {
	return p.conn.WriteMessage(messageType, data)
}

func (p *noProxy) Close() error {
	return p.conn.Close()
}

type WithoutVPNProxy struct{}

func (*WithoutVPNProxy) AttachVPN(conn *websocket.Conn) AppConn {
	return &noProxy{conn: conn}
}
