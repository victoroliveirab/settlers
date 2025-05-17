package utils

import (
	"encoding/json"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func ParseJsonPayload[T any](message *types.WebSocketClientRequest) (*T, error) {
	var payload T
	if err := json.Unmarshal(message.Payload, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
