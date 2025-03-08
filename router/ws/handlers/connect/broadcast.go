package connect

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func buildPlayerConnectedPreMatchBroadcast(player *entities.GamePlayer) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "room.player-connect",
		Payload: map[string]interface{}{
			"player": player.Username,
		},
	}
}

func buildPlayerReconnectedPreMatchBroadcast(player *entities.GamePlayer) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "room.player-reconnect",
		Payload: map[string]interface{}{
			"player": player.Username,
		},
	}
}
