package matchsetup

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendBuildSetupSettlementRequest(player *entities.GamePlayer) error {
	vertices, _ := player.Room.Game.AvailableVertices(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "setup.build-settlement",
		Payload: map[string]interface{}{
			"vertices": vertices,
		},
	})
}

func SendBuildSetupRoadRequest(player *entities.GamePlayer) error {
	edges, _ := player.Room.Game.AvailableEdges(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "setup.build-road",
		Payload: map[string]interface{}{
			"edges": edges,
		},
	})
}

func sendBuildSetupRoadRequest(player *entities.GamePlayer) error {
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "setup.build-road",
	})
}
