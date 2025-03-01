package match

import (
	"strconv"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func SendBuildSettlementRequest(conn *types.WebSocketConnection, userID int64, room *entities.Room) error {
	playerID := strconv.FormatInt(userID, 10)
	vertices, _ := room.Game.AvailableVertices(playerID)
	return utils.WriteJson(conn, userID, &types.WebSocketMessage{
		Type: "setup.settlement",
		Payload: map[string]interface{}{
			"vertices": vertices,
		},
	})
}
