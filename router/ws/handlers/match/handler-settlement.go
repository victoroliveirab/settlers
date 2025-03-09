package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

type newSettlementPayload struct {
	vertexID int
}

func parseNewSettlementPayload(payload map[string]interface{}) (*newSettlementPayload, error) {
	vertexID, ok := payload["vertex"].(float64)
	if !ok {
		err := fmt.Errorf("malformed data: vertex")
		return nil, err
	}

	return &newSettlementPayload{
		vertexID: int(vertexID),
	}, nil
}

func handleNewSettlement(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	payload, err := parseNewSettlementPayload(message.Payload)
	if err != nil {
		wsErr := sendNewSettlementError(player.Connection, player.ID, err)
		return true, wsErr
	}

	vertexID := payload.vertexID
	room := player.Room
	game := room.Game
	err = game.BuildSettlement(player.Username, vertexID)
	if err != nil {
		wsErr := sendNewSettlementError(player.Connection, player.ID, err)
		return true, wsErr
	}

	logs := []string{fmt.Sprintf("%s just built a road.", player.Username)}
	err = sendNewSettlementSuccess(player, vertexID, logs)
	room.EnqueueBroadcastMessage(buildNewSettlementBroadcast(player.Username, vertexID, logs), []int64{player.ID}, nil)
	return true, err

}
