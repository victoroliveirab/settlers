package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "room.update-param":
		if player.Username != player.Room.Owner {
			err := fmt.Errorf("Cannot update room param: not room owner")
			wsErr := sendUpdateParamError(player, err)
			return true, wsErr
		}
		payload, err := parseUpdateParamPayload(message.Payload)
		if err != nil {
			wsErr := sendUpdateParamError(player, err)
			return true, wsErr
		}

		room := player.Room
		err = room.UpdateParam(payload.Key, payload.Value)
		if err != nil {
			wsErr := sendUpdateParamError(player, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{}, nil)
		return true, nil
	case "room.toggle-ready":
		payload, err := parsePlayerReadyState(message.Payload)
		if err != nil {
			wsErr := sendToggleReadyRequestError(player, err)
			return true, wsErr
		}

		room := player.Room
		err = room.TogglePlayerReadyState(player.ID, payload.Ready)
		if err != nil {
			wsErr := sendToggleReadyRequestError(player, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{}, nil)
		return true, nil
	case "room.start-game":
		if player.Username != player.Room.Owner {
			err := fmt.Errorf("Cannot start game: not room owner")
			wsErr := sendStartGameRequestError(player, err)
			return true, wsErr
		}
		room := player.Room

		err := StartMatch(room)
		if err != nil {
			wsErr := sendStartGameRequestError(player, err)
			return true, wsErr
		}
		return true, nil
	default:
		return false, nil
	}
}
