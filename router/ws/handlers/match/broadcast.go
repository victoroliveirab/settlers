package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func buildPlayerReconnectedBroadcast(player *entities.GamePlayer) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.player-reconnect",
		Payload: map[string]interface{}{
			"player": player.Username,
			"bot":    false,
		},
	}
}

func buildSettlementSetupBuildSuccessBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:    "room.new-update",
		Payload: room.ToMapInterface(),
	}
}
