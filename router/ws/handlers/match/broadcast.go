package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildPlayerRoundBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.player-round",
		Payload: map[string]interface{}{
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
		},
	}
}

func buildMoveRobberDueTo7Broadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type: "game.move-robber-request",
		Payload: map[string]interface{}{
			"availableTiles":     room.Game.UnblockedTiles(),
			"currentRoundPlayer": room.Game.CurrentRoundPlayer().ID,
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

func buildDiscardCardsSuccessBroadcast(room *entities.Room, logs []string) *types.WebSocketMessage {
	quantityByPlayers := make(map[string]int)
	for _, participant := range room.Participants {
		username := participant.Player.Username
		quantityByPlayers[username] = room.Game.DiscardAmountByPlayer(username)
	}
	return &types.WebSocketMessage{
		Type: "game.discard-cards.success",
		Payload: map[string]interface{}{
			"quantityByPlayers": quantityByPlayers,
			"logs":              logs,
		},
	}

}
