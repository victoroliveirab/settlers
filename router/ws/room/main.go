package room

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/state"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

var color [4]string = [4]string{"green", "orange", "blue", "black"}

func addPlayerToGameRoom(roomID string, player *types.GamePlayer) error {
	room, exists := state.RoomByID[roomID]
	if !exists {
		err := fmt.Errorf("Cannot join room #%s: room doesn't exist", roomID)
		return err
	}

	for _, spot := range room.Participants {
		if spot.Player != nil {
			if spot.Player.ID == player.ID {
				return nil
			}
		}
	}

	for i, spot := range room.Participants {
		if spot.Player == nil {
			player.Color = color[i]
			room.Participants[i] = types.RoomEntry{
				Player: player,
				Ready:  false,
				Bot:    false,
			}
			return nil
		}
	}

	err := fmt.Errorf("Cannot join room #%s: room full", roomID)
	return err
}

func updatePlayerReadyness(roomID string, player *types.GamePlayer, ready bool) error {
	room, exists := state.RoomByID[roomID]
	if !exists {
		err := fmt.Errorf("Cannot change ready status in room #%s: room doesn't exist", roomID)
		return err
	}

	for i, entry := range room.Participants {
		if entry.Player.ID == player.ID {
			room.Participants[i].Ready = true
			return nil
		}
	}

	err := fmt.Errorf("Cannot change ready status in room #%s: user is not part of this room", roomID)
	return err
}

func HandleMessage(conn *types.WebSocketConnection, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "room.join":
		payload, err := readRoomPayload(message.Payload)
		if err != nil {
			wsErr := sendRoomJoinRequestError(conn, err)
			return true, wsErr
		}

		player := state.PlayerByConnection[conn]
		err = addPlayerToGameRoom(payload.RoomID, player)
		if err != nil {
			wsErr := sendRoomJoinRequestError(conn, err)
			return true, wsErr
		}

		err = sendRoomJoinRequestSuccess(conn, state.RoomByID[payload.RoomID])
		return true, err
	case "room.ready":
		payload, err := readRoomPayload(message.Payload)
		if err != nil {
			return true, err
		}

		player := state.PlayerByConnection[conn]
		err = updatePlayerReadyness(payload.RoomID, player, true)
		if err != nil {
			wsErr := sendRoomReadyRequestError(conn, err)
			return true, wsErr
		}

		err = sendRoomReadyRequestSuccess(conn, state.RoomByID[payload.RoomID])
		return true, err
	case "room.not-ready":
		payload, err := readRoomPayload(message.Payload)
		if err != nil {
			return true, err
		}

		player := state.PlayerByConnection[conn]
		err = updatePlayerReadyness(payload.RoomID, player, false)
		if err != nil {
			wsErr := sendRoomNotReadyRequestError(conn, err)
			return true, wsErr
		}

		err = sendRoomNotReadyRequestSuccess(conn, state.RoomByID[payload.RoomID])
		return true, err
	case "room.start":
		return true, nil
	default:
		return false, nil
	}
}
