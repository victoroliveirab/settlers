package utils

import (
	"encoding/json"
	"strings"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ReadJson(conn *types.WebSocketConnection, userID int64) (*types.WebSocketMessage, error) {
	var parsedMessage types.WebSocketMessage
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()
	m, message, err := conn.Instance.ReadMessage()
	if err != nil {
		logger.LogError(userID, "conn.ReadMessage", m, err)
		return nil, err
	}
	err = json.Unmarshal(message, &parsedMessage)
	if err != nil {
		logger.LogError(userID, "json.Unmarshal", -1, err)
		return nil, err
	}
	logger.LogWSMessage("incoming", userID, parsedMessage.Type, parsedMessage.Payload)
	return &parsedMessage, nil
}

func WriteJson(conn *types.WebSocketConnection, userID int64, message *types.WebSocketMessage) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()
	err := conn.Instance.WriteJSON(message)
	if err != nil {
		logger.LogError(userID, strings.Join([]string{"conn.WriteJSON", message.Type}, "."), -1, err)
		return err
	}
	logger.LogWSMessage("outgoing", userID, message.Type, message.Payload)
	return nil
}
