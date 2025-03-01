package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func TryHandle(room *entities.Room, player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "room.join":
		payload, err := parseRoomJoinPayload(message.Payload)
		if err != nil {
			wsErr := sendRoomJoinRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}

		if payload.RoomID != room.ID {
			wsErr := sendRoomJoinRequestError(player.Connection, player.ID, fmt.Errorf("Wrong roomID: %s", payload.RoomID))
			return true, wsErr
		}

		err = room.AddPlayer(player)
		if err != nil {
			wsErr := sendRoomJoinRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}

		err = sendRoomJoinRequestSuccess(player.Connection, player.ID, room)
		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{player.ID}, nil)
		return true, err
	case "room.toggle-ready":
		payload, err := parsePlayerReadyState(message.Payload)
		if err != nil {
			wsErr := sendToggleReadyRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}

		err = room.TogglePlayerReadyState(player.ID, payload.Ready)
		if err != nil {
			wsErr := sendToggleReadyRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}

		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{}, nil)
		return true, nil
	case "room.start-game":
		if player.ID != room.OwnerID {
			err := fmt.Errorf("Cannot start game: not room owner")
			wsErr := sendStartGameRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}

		err := StartMatch(room)
		if err != nil {
			wsErr := sendStartGameRequestError(player.Connection, player.ID, err)
			return true, wsErr
		}
		return true, nil
	default:
		return false, nil
	}
}
