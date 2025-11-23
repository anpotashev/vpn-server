package vpnproxy

import (
	"encoding/binary"
	"net"

	"github.com/anpotashev/vpn-server/internal/ipallocator"
	"github.com/gorilla/websocket"
)

// Описывает одного WS клиента
type client struct {
	conn             *websocket.Conn
	ip               *net.IP
	writeToInterface func([]byte) error
	ipAllocator      ipallocator.IPAllocator
	mtu              uint16
	beforeClose      func()
}

func (a *client) Read() (messageType int, data []byte, err error) {
	isWsMessage := true
	for isWsMessage {
		messageType, data, err = a.conn.ReadMessage()
		if err != nil {
			return messageType, data, err
		}
		isWsMessage, err = a.processWSMessage(messageType, data)
		if err != nil {
			return messageType, data, err
		}
	}
	return messageType, data, err
}

func (a *client) Write(messageType int, data []byte) error {
	return a.conn.WriteMessage(messageType, data)
}

func (a *client) Close() error {
	if a.ip != nil {
		a.ipAllocator.ReleaseIP(*a.ip)
	}
	a.beforeClose()
	return a.conn.Close()
}

func (a *client) processWSMessage(messageType int, data []byte) (bool, error) {
	if a.ip == nil {
		connectionSuccess, err := a.tryProcessMsgAsConnectCommand(messageType, data)
		if err != nil {
			return false, err
		}
		return connectionSuccess, nil
	}
	if messageType != websocket.BinaryMessage {
		return false, nil
	}
	if src := parseIPv4Src(data); src == nil || !src.Equal(*a.ip) {
		return false, nil
	}
	if err := a.writeToInterface(data); err != nil {
		return false, err
	}
	return true, nil
}

func (a *client) tryProcessMsgAsConnectCommand(messageType int, data []byte) (bool, error) {
	if messageType != websocket.TextMessage || string(data) != "start" {
		return false, nil
	}
	ip, err := a.ipAllocator.AllocateIP()
	if err != nil {
		return true, nil
	}
	a.ip = &ip
	dataToWrite := a.getBytesToSendIpInfo(ip)
	if err = a.conn.WriteMessage(websocket.BinaryMessage, dataToWrite); err != nil {
		return true, err
	}
	err = a.conn.WriteMessage(websocket.BinaryMessage, dataToWrite)
	if err != nil {
		a.ipAllocator.ReleaseIP(ip)
		a.ip = nil
		return true, err
	}
	return true, nil
}

func (a *client) getBytesToSendIpInfo(ip net.IP) []byte {
	packageSizeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(packageSizeBytes, a.mtu)
	// собираем и отправляем клиенту мета-пакет с информацией для клиента:
	dataToWrite := make([]byte, 0, 2+4+4+4+4)
	// 2 байта - макс размер пакета,
	dataToWrite = append(dataToWrite, packageSizeBytes...)
	// 4 байта - IP
	dataToWrite = append(dataToWrite, ip...)
	// 4 байта - шлюз
	dataToWrite = append(dataToWrite, a.ipAllocator.Gateway()...)
	// 4 байта - маска подсети
	dataToWrite = append(dataToWrite, a.ipAllocator.IPMask()...)
	return dataToWrite
}
