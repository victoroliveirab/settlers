package postmatch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildPostMatchMessage(room *entities.Room) *types.WebSocketServerResponse {
	game := room.Game
	return &types.WebSocketServerResponse{
		Type: "over.data",
		Payload: postMatchDataResponsePayload{
			Points:        game.Points(),
			RoomStatus:    room.Status,
			RoundsPlayed:  game.Round() + 1,
			Statistics:    game.GetStatistics(),
			StartDatetime: room.StartDatetime,
			EndDatetime:   room.EndDatetime,
		},
	}
}
