package room

import (
	"github.com/victoroliveirab/settlers/router/ws/state"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func roomToMapInterface(room *types.Room) map[string]interface{} {
	return map[string]interface{}{
		"id":           room.ID,
		"capacity":     room.Capacity,
		"map":          room.MapName,
		"participants": room.Participants,
	}
}

func sendRoomJoinRequestSuccess(conn *types.WebSocketConnection, room *types.Room) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type:    "room.join.success",
		Payload: roomToMapInterface(room),
	})
}

func sendRoomJoinRequestError(conn *types.WebSocketConnection, err error) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.join.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendRoomReadyRequestSuccess(conn *types.WebSocketConnection, room *types.Room) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type:    "room.ready.success",
		Payload: roomToMapInterface(room),
	})
}

func sendRoomReadyRequestError(conn *types.WebSocketConnection, err error) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.ready.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendRoomNotReadyRequestSuccess(conn *types.WebSocketConnection, room *types.Room) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type:    "room.not-ready.success",
		Payload: roomToMapInterface(room),
	})
}

func sendRoomNotReadyRequestError(conn *types.WebSocketConnection, err error) error {
	userID := state.PlayerByConnection[conn].ID
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.not-ready.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
