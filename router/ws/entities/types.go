package entities

import (
	"sync"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type Lobby struct {
	rooms        map[string]*Room
	roomByPlayer map[int64]*Room
	sync.Mutex
}

type RoomEntry struct {
	Player *GamePlayer `json:"player"`
	Ready  bool        `json:"ready"`
	Bot    bool        `json:"bot"`
}

type IncomingMessage struct {
	Room    *Room
	Player  *GamePlayer
	Message *types.WebSocketMessage
}

type BroadcastMessage struct {
	ExcludedIDs []int64
	Message     *types.WebSocketMessage
	OnSend      func()
}

type RoomIncomingMessageHandler func(player *GamePlayer, message *types.WebSocketMessage) (bool, error)

type Room struct {
	ID                string                       `json:"id"`
	Capacity          int                          `json:"capacity"`
	Game              *core.GameState              `json:"-"`
	MapName           string                       `json:"map"`
	Participants      []RoomEntry                  `json:"participants"`
	Private           bool                         `json:"private"`
	Owner             string                       `json:"owner"`
	incomingMsgQueue  chan IncomingMessage         `json:"-"`
	broadcastMsgQueue chan BroadcastMessage        `json:"-"`
	handlers          []RoomIncomingMessageHandler `json:"-"`
	onDestroy         func(room *Room)             `json:"-"`
	sync.Mutex
}

type GamePlayer struct {
	ID          int64                      `json:"-"`
	Username    string                     `json:"username"`
	Connection  *types.WebSocketConnection `json:"-"`
	Color       string                     `json:"color"`
	Room        *Room                      `json:"-"`
	OnDisconect func(player *GamePlayer)   `json:"-"`
}
