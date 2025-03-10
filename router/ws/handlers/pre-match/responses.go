package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func sendUpdateParamError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.update-param.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendToggleReadyRequestError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.toggle-ready.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendStartGameRequestError(player *entities.GamePlayer, err error) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "room.start-game.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
