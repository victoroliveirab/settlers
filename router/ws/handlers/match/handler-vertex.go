package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type vertexClickRequestPayload struct {
	VertexID int `json:"vertex"`
}

func handleVertexClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[vertexClickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	vertexID := payload.VertexID
	room := player.Room
	game := room.Game
	logs := make([]string, 0)

	_, exists := game.AllSettlements()[vertexID]
	if !exists {
		err := game.BuildSettlement(player.Username, vertexID)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}
		logs = append(logs, fmt.Sprintf("%s built a new settlement.", player.Username))
	} else {
		err := game.BuildCity(player.Username, vertexID)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}
		logs = append(logs, fmt.Sprintf("%s built a new city.", player.Username))
	}

	room.EnqueueBulkUpdate(
		UpdateMapState,
		UpdateEdgeState,
		UpdateVertexState,
		UpdatePlayerHand,
		UpdateBuyDevelopmentCard,
		UpdateLogs(logs),
	)
	return true, nil
}
