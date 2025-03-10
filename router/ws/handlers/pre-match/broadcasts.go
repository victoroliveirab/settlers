package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildRoomStateUpdateBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:    "room.new-update",
		Payload: room.ToMapInterface(),
	}
}

func buildStartGameBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.start",
		Payload: map[string]interface{}{
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
			"map":                room.Game.Map(),
			"players":            room.Game.Players(),
			"resourceCount":      room.Game.NumberOfResourcesByPlayer(), // not necessary, but maybe we want to give players initial resources sometime
			"logs":               logs,
		},
	}
}
