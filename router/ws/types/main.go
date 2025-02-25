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

type GamePlayer struct {
	ID         int64                `json:"id"`
	Username   string               `json:"username"`
	Connection *WebSocketConnection `json:"-"`
	Color      string               `json:"color"`
	Room       string               `json:"roomID"`
}

type RoomEntry struct {
	Player *GamePlayer `json:"player"`
	Ready  bool        `json:"ready"`
	Bot    bool        `json:"bot"`
}

type Room struct {
	ID           string      `json:"roomID"`
	Capacity     int         `json:"capacity"`
	MapName      string      `json:"map"`
	Participants []RoomEntry `json:"participants"`
}
