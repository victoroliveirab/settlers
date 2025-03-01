package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendHydratePlayer(player *entities.GamePlayer) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "hydrate",
		Payload: map[string]interface{}{
			"state": generateGameStateDump(player),
		},
	})
}

func SendReconnectPlayerError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.reconnect.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendSettlementSetupBuildError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "settlement.setup-build.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
