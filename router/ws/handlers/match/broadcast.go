package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildLogsBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.logs",
		Payload: map[string]interface{}{
			"logs": logs,
		},
	}
}

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

func buildMoveRobberBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.move-robber-request",
		Payload: map[string]interface{}{
			"availableTiles":     room.Game.UnblockedTiles(),
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
			"logs":               logs,
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

func buildNewSettlementBroadcast(builderID string, vertexID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.new-settlement.broadcast",
		Payload: map[string]interface{}{
			"settlement": map[string]interface{}{
				"id":    vertexID,
				"owner": builderID,
			},
			"logs": logs,
		},
	}
}

// REFACTOR: rebuild this to send all the blocked tiles instead (in case multiple robbers)
func buildRobberMovedBroadcast(tileID int, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.move-robber.broadcast",
		Payload: map[string]interface{}{
			"tile": tileID,
		},
	}
}

func buildRobFinishedBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:    "game.rob-finished",
		Payload: map[string]interface{}{},
	}
}
