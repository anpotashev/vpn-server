package ipallocator

import (
	"errors"
	"net"
	"sync"
)

type IPAllocator interface {
	AllocateIP() (net.IP, error)
	ReleaseIP(ip net.IP)
	Gateway() net.IP
	IPMask() net.IPMask
}

var ErrNoIpAvailable = errors.New("no IP available")

type allocator struct {
	mu           sync.Mutex
	ipTemplate   []byte
	allocatedMap map[byte]bool
}

func NewIP4Allocator(ipTemplate []byte) IPAllocator {
	return &allocator{
		allocatedMap: make(map[byte]bool),
		ipTemplate:   ipTemplate,
	}
}

func (a *allocator) AllocateIP() (net.IP, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 2; i < 255; i++ {
		if used := a.allocatedMap[byte(i)]; used {
			continue
		}
		a.allocatedMap[byte(i)] = true
		return append(a.ipTemplate, byte(i)), nil
	}
	return nil, ErrNoIpAvailable
}

func (a *allocator) ReleaseIP(ip net.IP) {
	delete(a.allocatedMap, ip[3])
}
func (a *allocator) Gateway() net.IP {
	return append(a.ipTemplate, byte(1))
}
func (a *allocator) IPMask() net.IPMask {
	return []byte{255, 255, 255, 0}
}
