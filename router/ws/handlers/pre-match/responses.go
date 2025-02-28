package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func sendRoomJoinRequestSuccess(conn *types.WebSocketConnection, userID int64, room *entities.Room) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type:    "room.join.success",
		Payload: room.ToMapInterface(),
	})
}

func sendRoomJoinRequestError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.join.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func sendToggleReadyRequestError(conn *types.WebSocketConnection, userID int64, err error) error {
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "room.toggle-ready.error",
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
}
