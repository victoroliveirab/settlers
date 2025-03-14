package matchsetup

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func UpdateCurrentRoundPlayerState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-round-player", room.Status)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: currentRoundPlayerStateUpdateResponsePayload{
			Player: game.CurrentRoundPlayer().ID,
		},
	}
}

func UpdateVertexState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-vertices", room.Status)
	availableVertices, err := game.AvailableVertices(username)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: verticesStateUpdateResponsePayload{
			AvailableVertices: availableVertices,
			Disabled:          err != nil,
		},
	}
}

func UpdateEdgeState(room *entities.Room, username string) *types.WebSocketServerResponse {
	game := room.Game
	messageType := fmt.Sprintf("%s.update-edges", room.Status)
	availableEdges, err := game.AvailableEdges(username)
	return &types.WebSocketServerResponse{
		Type: types.ResponseType(messageType),
		Payload: edgesStateUpdateResponsePayload{
			AvailableEdges: availableEdges,
			Disabled:       err != nil,
		},
	}
}

func UpdateLogs(logs []string) func(room *entities.Room, username string) *types.WebSocketServerResponse {
	return func(room *entities.Room, username string) *types.WebSocketServerResponse {
		messageType := fmt.Sprintf("%s.update-edges", room.Status)
		return &types.WebSocketServerResponse{
			Type:    types.ResponseType(messageType),
			Payload: logs,
		}
	}
}
