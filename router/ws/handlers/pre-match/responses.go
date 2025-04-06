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

func buildStartMatch(room *entities.Room) *types.WebSocketServerResponse {
	game := room.Game
	responsePayload := roomStartMatchPayload{
		Map:           game.Map(),
		MapName:       "base4",
		Players:       game.Players(),
		ResourceCount: game.NumberOfResourcesByPlayer(),
		Logs:          []string{},
	}
	msg := &types.WebSocketServerResponse{
		Type:    types.ResponseType("room.start-game.success"),
		Payload: responsePayload,
	}
	return msg
}
