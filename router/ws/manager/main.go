package manager

import (
	"github.com/gorilla/websocket"

	core "github.com/victoroliveirab/settlers/core/state"
)

type GamePlayer struct {
	ID       int64
	Username string
	Color    string
}

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
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

type WSServerMaps struct {
	PlayerMap          map[int64]*GamePlayer
	ConnectionByPlayer map[int64]*websocket.Conn
	GameByPlayer       map[int64]*core.GameState
	UsersIDsByGame     map[*core.GameState][]int64
	RoomByID           map[string]Room
}
type Data struct {
	Connection *websocket.Conn
	Internal   WSServerMaps
	Message    Message
	UserID     string
}

var maps WSServerMaps = WSServerMaps{
	PlayerMap:          map[int64]*GamePlayer{},
	ConnectionByPlayer: map[int64]*websocket.Conn{},
	GameByPlayer:       map[int64]*core.GameState{},
	UsersIDsByGame:     map[*core.GameState][]int64{},
	RoomByID:           map[string]Room{},
}
