package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	switch message.Type {
	case "match.dice-roll":
		return handleDiceRoll(player, message)
	case "match.vertex-click":
		return handleVertexClick(player, message)
	case "match.edge-click":
		return handleEdgeClick(player, message)
	case "match.pass-click":
		return handleEndRound(player, message)
	case "match.create-trade-offer":
		return handleCreateTradeOffer(player, message)
	case "match.create-counter-trade-offer":
		return handleCreateCounterTradeOffer(player, message)
	case "match.accept-trade-offer":
		return handleAcceptTradeOffer(player, message)
	case "match.reject-trade-offer":
		return handleRejectTradeOffer(player, message)
	case "match.cancel-trade-offer":
		return handleCancelTradeOffer(player, message)
	case "match.finalize-trade-offer":
		return handleFinalizeTradeOffer(player, message)
	// case "game.discard-cards":
	// 	return handleDiscardCards(player, message)
	// case "game.move-robber":
	// 	return handleMoveRobber(player, message)
	case "match.end-round":
		return handleEndRound(player, message)
	default:
		return false, nil
	}
}
