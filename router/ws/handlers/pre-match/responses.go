package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildRoomMessage(room *entities.Room, messageType string) *types.WebSocketServerResponse {
	responsePayload := roomUpdateResponsePayload{
		Room:       room,
		RoomParams: room.Params(),
	}
	msg := &types.WebSocketServerResponse{
		Type:    types.ResponseType(messageType),
		Payload: responsePayload,
	}
	return msg
}
