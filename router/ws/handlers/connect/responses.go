package connect

import (
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func sendReconnectPlayerError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.reconnect.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
