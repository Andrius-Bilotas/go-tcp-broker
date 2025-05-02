package cache

import (
	"net"
	"slices"
	"sync"
)

type Cache interface {
	GetConnections(key ConnectionType) []net.Conn
	AddConnection(key ConnectionType, value net.Conn) []net.Conn
	RemoveConnection(key ConnectionType, value net.Conn) []net.Conn
}

type ConnectionCache struct {
	data map[ConnectionType][]net.Conn
	mu   sync.Mutex
}

type ConnectionType string

const (
	Publisher  ConnectionType = "publishers"
	Subscriber ConnectionType = "subscribers"
)

func NewConnectionCache() *ConnectionCache {
	return &ConnectionCache{
		data: make(map[ConnectionType][]net.Conn),
	}
}

func (c *ConnectionCache) GetConnections(key ConnectionType) []net.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	val := c.data[key]
	return val
}

func (c *ConnectionCache) AddConnection(key ConnectionType, value net.Conn) []net.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = append(c.data[key], value)
	return c.data[key]
}

func (c *ConnectionCache) RemoveConnection(key ConnectionType, value net.Conn) []net.Conn {
	c.mu.Lock()
	defer c.mu.Unlock()
	connections := c.data[key]

	for i, conn := range connections {
		if value == conn {
			connections = slices.Delete(connections, i, i+1)
			break
		}
	}
	c.data[key] = connections

	return c.data[key]
}
