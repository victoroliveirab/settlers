package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type edgeClickRequestPayload struct {
	EdgeID int `json:"edge"`
}

func handleEdgeClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[edgeClickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	edgeID := payload.EdgeID
	room := player.Room
	game := room.Game

	err = game.BuildRoad(player.Username, edgeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	logs := make([]string, 1)
	logs[0] = fmt.Sprintf("%s has built a new road.", player.Username)
	if game.RoundType() == core.FirstRound {
		room.ProgressStatus()
		logs = append(logs, "Setup phase is over.", "Match starting. Good luck to everyone!")
		room.EnqueueBulkUpdate(
			UpdateMapState,
			UpdateCurrentRoundPlayerState,
			UpdateVertexState,
			UpdateEdgeState,
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdateDiceState,
			UpdateBuyDevelopmentCard,
			UpdateLongestRoadSize,
			UpdateLogs(logs),
		)
	} else {
		room.EnqueueBulkUpdate(
			UpdateMapState,
			UpdateCurrentRoundPlayerState,
			UpdateVertexState,
			UpdateEdgeState,
			UpdatePlayerHand,
			UpdateBuyDevelopmentCard,
			UpdateLongestRoadSize,
			UpdatePoints,
			UpdateLogs(logs),
		)
	}
	return true, nil
}
