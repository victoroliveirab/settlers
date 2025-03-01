package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(room *entities.Room, player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "settlement.setup-build":
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

		room.EnqueueBroadcastMessage(buildSettlementSetupBuildSuccessBroadcast(room), []int64{}, func() {
			// roundPlayer := room.Game.CurrentRoundPlayer()
			// roundPlayer.ID
		})
		return true, nil
	default:
		return false, nil
	}
}
