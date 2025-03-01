package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "room.join":
		payload, err := parseRoomJoinPayload(message.Payload)
		if err != nil {
			wsErr := sendRoomJoinRequestError(player, err)
			return true, wsErr
		}

		if payload.RoomID != player.Room.ID {
			wsErr := sendRoomJoinRequestError(player, fmt.Errorf("Wrong roomID: %s", payload.RoomID))
			return true, wsErr
		}
		room := player.Room

		err = room.AddPlayer(player)
		if err != nil {
			wsErr := sendRoomJoinRequestError(player, err)
			return true, wsErr
		}

		err = sendRoomJoinRequestSuccess(player)
		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{player.ID}, nil)
		return true, err
	case "room.toggle-ready":
		// TODO: require room id in request
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
