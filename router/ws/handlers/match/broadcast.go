package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildPlayerRoundOpponentsBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.player-round-changed",
		Payload: map[string]interface{}{
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
			"roundType":          room.Game.RoundType(),
			"round":              room.Game.Round(),
		},
	}
}

func buildMoveRobberDueTo7Broadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.move-robber-request",
		Payload: map[string]interface{}{
			"availableTiles": room.Game.UnblockedTiles(),
		},
	}
}

func buildDiscardCardsBroadcast(room *entities.Room) *types.WebSocketMessage {
	quantityByPlayers := make(map[string]int)
	for _, participant := range room.Participants {
		username := participant.Player.Username
		quantityByPlayers[username] = room.Game.DiscardAmountByPlayer(username)
	}
	return &types.WebSocketMessage{
		Type: "game.discard-cards-request",
		Payload: map[string]interface{}{
			"quantityByPlayers": quantityByPlayers,
		},
	}
}

func buildDiscardedCardsBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	quantityByPlayers := make(map[string]int)
	for _, participant := range room.Participants {
		username := participant.Player.Username
		quantityByPlayers[username] = room.Game.DiscardAmountByPlayer(username)
	}
	return &types.WebSocketMessage{
		Type: "game.discarded-cards",
		Payload: map[string]interface{}{
			"resourceCount":     room.Game.NumberOfResourcesByPlayer(),
			"quantityByPlayers": quantityByPlayers,
			"logs":              logs,
		},
	}
}

func buildNewRoadBroadcast(builderID string, edgeID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.new-road.broadcast",
		Payload: map[string]interface{}{
			"road": map[string]interface{}{
				"id":    edgeID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}
