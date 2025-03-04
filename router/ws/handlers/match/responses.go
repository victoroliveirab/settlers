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

func sendDiscardCardsError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "game.discard-cards.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
