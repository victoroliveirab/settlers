package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type WebSocketConnection struct {
	Instance *websocket.Conn
	Mutex    sync.Mutex
}
