package connections

import (
	"sync"

	"github.com/gorilla/websocket"
)

// connectionMap stores user ID -> WebSocket connection
var connectionMap = sync.Map{}

// AddConnection adds a WebSocket connection for a user
func AddConnection(userID string, conn *websocket.Conn) {
	connectionMap.Store(userID, conn)
}

// RemoveConnection removes a WebSocket connection for a user
func RemoveConnection(userID string) {
	connectionMap.Delete(userID)
}

// GetConnection retrieves the WebSocket connection for a user
func GetConnection(userID string) (*websocket.Conn, bool) {
	conn, ok := connectionMap.Load(userID)
	if !ok {
		return nil, false
	}
	return conn.(*websocket.Conn), true
}
