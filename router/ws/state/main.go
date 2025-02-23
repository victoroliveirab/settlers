package state

import (
	"github.com/gorilla/websocket"
	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

var PlayerMap = map[int64]*types.GamePlayer{}
var ConnectionByPlayer = map[int64]*websocket.Conn{}
var GameByPlayer = map[int64]*core.GameState{}
var UsersIDsByGame = map[*core.GameState][]int64{}
var RoomByID = map[string]types.Room{}
