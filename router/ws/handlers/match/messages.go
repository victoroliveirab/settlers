package match

import (
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendBuildSettlementRequest(player *entities.GamePlayer) error {
	vertices, _ := player.Room.Game.AvailableVertices(player.Username)
	return utils.WriteJson(player.Connection, player.ID, &types.WebSocketMessage{
		Type: "setup.settlement",
		Payload: map[string]interface{}{
			"vertices": vertices,
		},
	})
}
