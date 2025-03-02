package prematch

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/reconnect"
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

		// ongoing match, player may be trying to reconnect
		if room.Game != nil {
			err := reconnect.TryReconnectPlayer(player)
			if err != nil {
				wsErr := reconnect.SendReconnectPlayerError(player, err)
				return true, wsErr
			}
			return true, nil
		}

		err = room.AddPlayer(player)
		if err != nil {
			wsErr := sendRoomJoinRequestError(player, err)
			return true, wsErr
		}

		err = sendRoomJoinRequestSuccess(player)
		room.EnqueueBroadcastMessage(buildRoomStateUpdateBroadcast(room), []int64{player.ID}, nil)
		return true, err
	case "room.toggle-ready":
		payload, err := parsePlayerReadyState(message.Payload)
		if err != nil {
			wsErr := sendToggleReadyRequestError(player, err)
			return true, wsErr
		}

		if payload.RoomID != player.Room.ID {
			wsErr := sendRoomJoinRequestError(player, fmt.Errorf("Cannot toggle ready in room#%s: not part of room", payload.RoomID))
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
