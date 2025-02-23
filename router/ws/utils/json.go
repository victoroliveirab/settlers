package utils

import (
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/logger"
	types "github.com/victoroliveirab/settlers/router/ws/manager"
)

func ReadJson(conn *websocket.Conn, userID string) (*types.Message, error) {
	var parsedMessage types.Message
	m, message, err := conn.ReadMessage()
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

func WriteJson(conn *websocket.Conn, userID string, message *types.Message) error {
	err := conn.WriteJSON(message)
	if err != nil {
		logger.LogError(userID, strings.Join([]string{"conn.WriteJSON", message.Type}, "."), -1, err)
		return err
	}
	logger.LogWSMessage("outgoing", userID, message.Type, message.Payload)
	return nil
}
