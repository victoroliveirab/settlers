package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "game.dice-roll":
		return handleDiceRoll(player, message)
	case "game.discard-cards":
		return handleDiscardCards(player, message)
	case "game.new-road":
		return handleNewRoad(player, message)
	case "game.new-settlement":
		return handleNewSettlement(player, message)
	case "game.end-round":
		return handleEndRound(player, message)
	default:
		return false, nil
	}
}
