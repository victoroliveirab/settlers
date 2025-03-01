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

func buildSettlementSetupBuildSuccessBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "setup.settlement-build.success",
		Payload: map[string]interface{}{
			"settlements": room.Game.AllSettlements(),
			"logs":        logs,
		},
	}
}
