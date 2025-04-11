package types

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type RequestType string
type ResponseType string

type WebSocketConnection struct {
	Instance *websocket.Conn
	Mutex    sync.Mutex
}

type WebSocketClientRequest struct {
	Type    RequestType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WebSocketServerResponse struct {
	Type    ResponseType `json:"type"`
	Payload interface{}  `json:"payload"`
}
