package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
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
	currentRoundType := game.RoundType()

	err = game.BuildRoad(player.Username, edgeID)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	logs := []string{fmt.Sprintf("%s has built a new road.", player.Username)}

	if room.Status == "setup" {
		handleEdgeClickSetupResponse(room, logs)
	} else {
		handleEdgeClickMatchResponse(room, currentRoundType, logs)
	}
	return true, nil
}

func handleEdgeClickSetupResponse(room *entities.Room, logs []string) {
	game := room.Game
	if game.RoundType() == round.FirstRound {
		room.ProgressStatus()
		logs = append(logs, "Setup phase is over.", "Match starting. Good luck to everyone!")
		room.StartSubRound(round.FirstRound)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateMapState,
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
		room.EndRound()
		room.StartRound()
		room.StartSubRound(game.RoundType())
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateMapState,
			UpdateVertexState,
			UpdateEdgeState,
			UpdatePlayerHand,
			UpdateBuyDevelopmentCard,
			UpdateLongestRoadSize,
			UpdatePoints,
			UpdateLogs(logs),
		)
	}
}

func handleEdgeClickMatchResponse(room *entities.Room, prevRoundType round.Type, logs []string) {
	if prevRoundType == round.BuildRoad1Development {
		room.StartSubRound(round.BuildRoad2Development)
	} else if prevRoundType == round.BuildRoad2Development {
		room.ResumeRound()
	}

	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateMapState,
		UpdateVertexState,
		UpdateEdgeState,
		UpdatePlayerHand,
		UpdateBuyDevelopmentCard,
		UpdateLongestRoadSize,
		UpdatePoints,
		UpdateLogs(logs),
	)
}
