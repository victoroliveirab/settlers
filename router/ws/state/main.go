package state

import (
	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

// Maps Connection with actual Player
var PlayerByConnection = map[*types.WebSocketConnection]*types.GamePlayer{}

// Maps game state with participant players
var UsersIDsByGame = map[*core.GameState][]int64{}

// Maps Room.ID with actual Room
var RoomByID = map[string]*types.Room{}

// Maps Room.ID with game state
var GameByRoom = map[string]*core.GameState{}
