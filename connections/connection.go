package connections

import (
	"fmt"
	"sync"

	"websocket-server/logging"
	"websocket-server/utils"

	"github.com/gorilla/websocket"
)

// connectionMap stores user connections in a thread-safe manner
var connectionMap = sync.Map{}

// AddConnection adds a WebSocket connection for a user
func AddConnection(userID string, token string, conn *websocket.Conn) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		logging.LogError("Invalid token", fmt.Errorf("invalid token"))
		return fmt.Errorf("invalid token")
	}
	connectionMap.Store(userID, conn)
	logging.LogEvent(fmt.Sprintf("Connection added for userID: %s", userID))
	return nil
}

// RemoveConnection removes a WebSocket connection for a user
func RemoveConnection(userID string, token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		logging.LogError("Invalid token", fmt.Errorf("invalid token"))
		return fmt.Errorf("invalid token")
	}
	connectionMap.Delete(userID)
	logging.LogEvent(fmt.Sprintf("Connection removed for userID: %s", userID))
	return nil
}

// GetConnection retrieves the WebSocket connection for a user
func GetConnection(userID string, token string) (*websocket.Conn, error) {
	// _, err := utils.ValidateToken(token)
	// if err != nil {
	// 	logging.LogError("Invalid token", err)
	// 	// return nil, err
	// }
	conn, ok := connectionMap.Load(userID)
	if !ok {
		logging.LogError("Connection not found", fmt.Errorf("connection not found for userID: %s", userID))
		return nil, fmt.Errorf("connection not found")
	}
	logging.LogEvent(fmt.Sprintf("Connection retrieved for userID: %s", userID))
	return conn.(*websocket.Conn), nil
}
