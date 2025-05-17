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

// REFACTOR: replace with separate calls for settlement and city
func handleVertexClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[vertexClickRequestPayload](message)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	vertexID := payload.VertexID
	room := player.Room
	game := room.Game
	logs := make([]string, 0)

	_, exists := game.GetAllSettlements()[vertexID]
	if !exists {
		err := game.BuildSettlement(player.Username, vertexID)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}
		logs = append(logs, fmt.Sprintf("%s built a new settlement.", player.Username))
	} else {
		err := game.BuildCity(player.Username, vertexID)
		if err != nil {
			wsErr := player.WriteJsonError(message.Type, err)
			return true, wsErr
		}
		logs = append(logs, fmt.Sprintf("%s built a new city.", player.Username))
	}

	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateMapState,
		UpdateEdgeState,
		UpdateVertexState,
		UpdatePlayerHand,
		UpdatePortsState,
		UpdateBuyDevelopmentCard,
		UpdatePoints,
		UpdatePortsState,
		UpdateLongestRoadSize,
		UpdateLogs(logs),
	)
	return true, nil
}
