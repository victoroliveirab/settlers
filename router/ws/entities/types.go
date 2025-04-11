package entities

import (
	"sync"

	"github.com/victoroliveirab/settlers/core"
	coreT "github.com/victoroliveirab/settlers/core/types"
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
	Player  *GamePlayer
	Message *types.WebSocketClientRequest
}

type OutgoingMessage struct {
	Callback   func()
	Message    *types.WebSocketServerResponse
	Recipients []string
}

type RoomIncomingMessageHandler func(player *GamePlayer, message *types.WebSocketClientRequest) (bool, error)

type RoomParamsMetaEntry struct {
	Description string `json:"description"`
	Key         string `json:"key"`
	Label       string `json:"label"`
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
	ID               string                       `json:"id"`
	Capacity         int                          `json:"capacity"`
	Game             *core.GameState              `json:"-"`
	MapName          string                       `json:"map"`
	params           RoomParams                   `json:"-"`
	Participants     []RoomEntry                  `json:"participants"`
	Private          bool                         `json:"private"`
	Owner            string                       `json:"owner"`
	Status           string                       `json:"status"`
	Colors           []coreT.PlayerColor                `json:"colors"`
	incomingMsgQueue chan IncomingMessage         `json:"-"`
	outgoingMsgQueue chan OutgoingMessage         `json:"-"`
	handlers         []RoomIncomingMessageHandler `json:"-"`
	onDestroy        func(room *Room)             `json:"-"`
	sync.Mutex
}

type GamePlayer struct {
	ID           int64                      `json:"-"`
	Username     string                     `json:"name"`
	Connection   *types.WebSocketConnection `json:"-"`
	Color        *coreT.PlayerColor               `json:"color"`
	Room         *Room                      `json:"-"`
	OnDisconnect func(player *GamePlayer)   `json:"-"`
}
