package socket

import (
	"context"
	"github.com/gofiber/websocket/v2"
	"sync"
	"wormhole/pkg/database"
)

type ConnectionManager struct {
	connections map[*websocket.Conn]uint64
	mutex       sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[*websocket.Conn]uint64),
	}
}

func (cm *ConnectionManager) AddConnection(con *websocket.Conn, accId uint64) error {
	if err := cm.SetOnline(accId, true); err != nil {
		return err
	}

	cm.mutex.Lock()
	cm.connections[con] = accId
	cm.mutex.Unlock()

	return nil
}

func (cm *ConnectionManager) RemoveConnection(con *websocket.Conn) error {
	cm.mutex.RLock()
	accId, ok := cm.connections[con]
	cm.mutex.RUnlock()

	if !ok {
		return nil
	}

	if err := cm.SetOnline(accId, false); err != nil {
		return err
	}

	cm.mutex.Lock()
	delete(cm.connections, con)
	cm.mutex.Unlock()

	return nil
}

func (cm *ConnectionManager) SetOnline(accId uint64, online bool) error {
	query := "UPDATE group_members SET online = $1 WHERE account_id = $2"
	_, err := database.Get().Exec(context.Background(), query, online, accId)
	return err
}

func (cm *ConnectionManager) SendMessage(groupID int64, message []byte) {

}
