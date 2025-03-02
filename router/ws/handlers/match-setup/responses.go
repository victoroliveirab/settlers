package matchsetup

import (
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func sendSettlementSetupBuildError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "setup.new-settlement.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendRoadSetupBuildError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "setup.new-road.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
