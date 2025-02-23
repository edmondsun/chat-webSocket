// model/room.go
package model

import "sync"

// Room holds connected clients in a chat room.
type Room struct {
	Name    string
	Clients map[string]*ClientConn // Map of client ID to connection info.
	Mutex   sync.RWMutex           // Protects Clients map.
}

// ClientConn holds minimal connection info for a client in a room.
type ClientConn struct {
	ID   string
	Conn *Client // Pointer to the Client.
}
