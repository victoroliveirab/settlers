package entities

import (
	"sync"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type Lobby struct {
	availableRooms []string
	rooms          map[string]*Room
	roomByPlayer   map[int64]*Room
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

type RoomParamsMetaEntry struct {
	Description string `json:"description"`
	Key         string `json:"key"`
	Priority    int    `json:"-"`
	Value       int    `json:"value"`
	Values      []int  `json:"values"`
}
type RoomParamsMeta map[string]RoomParamsMetaEntry

type RoomParams struct {
	Meta   RoomParamsMeta `json:"meta"`
	Values map[string]int `json:"values"`
}

type Room struct {
	ID                string                       `json:"id"`
	Capacity          int                          `json:"capacity"`
	Game              *core.GameState              `json:"-"`
	MapName           string                       `json:"map"`
	params            RoomParams                   `json:"-"`
	Participants      []RoomEntry                  `json:"participants"`
	Private           bool                         `json:"private"`
	Owner             string                       `json:"owner"`
	Status            string                       `json:"status"`
	incomingMsgQueue  chan IncomingMessage         `json:"-"`
	broadcastMsgQueue chan BroadcastMessage        `json:"-"`
	handlers          []RoomIncomingMessageHandler `json:"-"`
	onDestroy         func(room *Room)             `json:"-"`
	sync.Mutex
}

type GamePlayer struct {
	ID           int64                      `json:"-"`
	Username     string                     `json:"name"`
	Connection   *types.WebSocketConnection `json:"-"`
	Color        string                     `json:"color"`
	Room         *Room                      `json:"-"`
	OnDisconnect func(player *GamePlayer)   `json:"-"`
}
