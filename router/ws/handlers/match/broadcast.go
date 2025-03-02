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

func buildSettlementSetupBuildSuccessBroadcast(builderID string, vertexID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "setup.settlement-build.success",
		Payload: map[string]interface{}{
			"settlement": map[string]interface{}{
				"id":    vertexID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}

func buildRoadSetupBuildSuccessBroadcast(builderID string, edgeID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "setup.road-build.success",
		Payload: map[string]interface{}{
			"road": map[string]interface{}{
				"id":    edgeID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}
