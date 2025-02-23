package types

import (
	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/core"
)

type WebSocketMessage[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}

type GamePlayer struct {
	ID       int64
	Username string
	Color    string
}

type RoomEntry struct {
	Connection *websocket.Conn
	UserID     int
	Color      string
	Bot        bool
}

type Room struct {
	ID       string
	Capacity int
	MapName  string
	Players  []RoomEntry
}

type WSState struct {
	PlayerMap          map[int64]*GamePlayer
	ConnectionByPlayer map[int64]*websocket.Conn
	GameByPlayer       map[int64]*core.GameState
	UsersIDsByGame     map[*core.GameState][]int64
	RoomByID           map[string]Room
}

type Data struct {
	Connection *websocket.Conn
	Internal   WSServerMaps
	Message    WebSocketMessage[any]
	UserID     string
}

var maps WSServerMaps = WSServerMaps{
	PlayerMap:          map[int64]*GamePlayer{},
	ConnectionByPlayer: map[int64]*websocket.Conn{},
	GameByPlayer:       map[int64]*core.GameState{},
	UsersIDsByGame:     map[*core.GameState][]int64{},
	RoomByID:           map[string]Room{},
}
