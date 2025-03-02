package reconnect

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendReconnectPlayerError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.reconnect.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
