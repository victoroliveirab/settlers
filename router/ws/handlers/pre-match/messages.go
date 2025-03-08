package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendConnectPlayerSuccess(player *entities.GamePlayer) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type:    "room.connect.success",
		Payload: player.Room.ToMapInterface(),
	})
}

func SendConnectPlayerError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.connect.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
