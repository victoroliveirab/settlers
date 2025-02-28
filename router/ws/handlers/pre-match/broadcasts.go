package prematch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func buildRoomStateUpdateBroadcast(room *entities.Room) *types.WebSocketMessage {
	return &types.WebSocketMessage{
		Type:    "room.new-update",
		Payload: room.ToMapInterface(),
	}
}
