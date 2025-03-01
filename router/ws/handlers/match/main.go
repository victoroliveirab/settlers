package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(room *entities.Room, player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "setup.new-settlement":
		payload, err := parseSettlementBuildPayload(message.Payload)
		if err != nil {
			wsErr := sendSettlementSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		vertexID := payload.vertexID
		err = room.Game.BuildSettlement(player.Username, vertexID)
		if err != nil {
			wsErr := sendSettlementSetupBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildSettlementSetupBuildSuccessBroadcast(room, []string{fmt.Sprintf("%s just built a settlement.", player.Username)}), []int64{}, func() {
			// TODO: send to player the request to build road
		})
		return true, nil
	default:
		return false, nil
	}
}
