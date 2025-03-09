package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func buildDiceRollSuccess(room *entities.Room, player *entities.GamePlayer, logs []string) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.dice-roll.success",
		Payload: map[string]interface{}{
			"dices":         room.Game.Dice(),
			"hand":          room.Game.ResourceHandByPlayer(player.Username),
			"logs":          logs,
			"resourceCount": room.Game.NumberOfResourcesByPlayer(),
			"roundType":     room.Game.RoundType(),
		},
	}
}

func sendDiceRollError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.dice-roll.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendEndRoundError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.end-round.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendDiscardCardsSuccess(room *entities.Room, player *entities.GamePlayer) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "game.discard-cards.success",
		Payload: map[string]interface{}{
			"hand":          room.Game.ResourceHandByPlayer(player.Username),
			"resourceCount": room.Game.NumberOfResourcesByPlayer(), // not necessary, but maybe we want to give players initial resources sometime
		},
	})
}

func sendDiscardCardsError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.discard-cards.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendNewRoadSuccess(player *entities.GamePlayer, edgeID int, logs []string) error {
	game := player.Room.Game
	availableEdges, _ := game.AvailableEdges(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "game.new-road.success",
		Payload: map[string]interface{}{
			"availableEdges": availableEdges,
			"hand":           game.ResourceHandByPlayer(player.Username),
			"road": map[string]interface{}{
				"id":    edgeID,
				"owner": player.Username,
			},
			"logs": logs,
		},
	})
}

func sendNewRoadError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.new-road.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendNewSettlementSuccess(player *entities.GamePlayer, vertexID int, logs []string) error {
	game := player.Room.Game
	availableVertices, _ := game.AvailableVertices(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "game.new-road.success",
		Payload: map[string]interface{}{
			"availableVertices": availableVertices,
			"hand":              game.ResourceHandByPlayer(player.Username),
			"settlement": map[string]interface{}{
				"id":    vertexID,
				"owner": player.Username,
			},
			"logs": logs,
		},
	})
}

func sendNewSettlementError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.new-settlement.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
