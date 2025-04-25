package postmatch

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func BuildPostMatchMessage(room *entities.Room) *types.WebSocketServerResponse {
	game := room.Game
	report := game.GetReport()

	return &types.WebSocketServerResponse{
		Type: "over.data",
		Payload: postMatchDataResponsePayload{
			Report:        report,
			RoomStatus:    room.Status,
			RoundsPlayed:  game.Round() + 1,
			StartDatetime: room.StartDatetime,
			EndDatetime:   room.EndDatetime,
		},
	}
}

func BuildPostMatchHydrateMessage(room *entities.Room) *types.WebSocketServerResponse {
	game := room.Game
	report := game.GetReport()

	return &types.WebSocketServerResponse{
		Type: "over.hydrate",
		Payload: postMatchHydrateResponsePayload{
			Report:        report,
			RoomName:      room.ID,
			RoomStatus:    room.Status,
			RoundsPlayed:  game.Round() + 1,
			StartDatetime: room.StartDatetime,
			EndDatetime:   room.EndDatetime,
			Players:       game.Players(),
			MapName:       game.MapName(),
		},
	}
}
