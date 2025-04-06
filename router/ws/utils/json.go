package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ReadJson(conn *types.WebSocketConnection, userID int64) (*types.WebSocketClientRequest, error) {
	var parsedMessage types.WebSocketClientRequest
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
	logger.LogWSMessage("incoming", userID, string(parsedMessage.Type), parsedMessage.Payload)
	return &parsedMessage, nil
}

func ParseJsonPayload[T any](message *types.WebSocketClientRequest) (*T, error) {
	var payload T
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func WriteJson(conn *types.WebSocketConnection, userID int64, message *types.WebSocketServerResponse) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()
	err := conn.Instance.WriteJSON(message)
	if err != nil {
		logger.LogError(userID, strings.Join([]string{"conn.WriteJSON", string(message.Type)}, "."), -1, err)
		return err
	}
	logger.LogWSMessage("outgoing", userID, string(message.Type), message.Payload)
	return nil
}

func WriteJsonError(conn *types.WebSocketConnection, userID int64, requestType types.RequestType, err error) error {
	message := &types.WebSocketServerResponse{
		Type: types.ResponseType(fmt.Sprintf("%s.error", requestType)),
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	}
	return WriteJson(conn, userID, message)
}
