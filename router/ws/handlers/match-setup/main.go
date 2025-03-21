package matchsetup

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/handlers/match"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func TryHandle(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	switch message.Type {
	case "setup.new-settlement":
		requestPayload, err := utils.ParseJsonPayload[newSettlementRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		vertexID := requestPayload.VertexID
		room := player.Room
		game := room.Game

		err = game.BuildSettlement(player.Username, vertexID)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		room.EnqueueBulkUpdate(
			match.UpdateMapState,
			match.UpdateEdgeState,
			match.UpdateLogs([]string{fmt.Sprintf("%s has built a new settlement.", player.Username)}),
		)
		return true, nil
	case "setup.new-road":
		requestPayload, err := utils.ParseJsonPayload[newRoadRequestPayload](message)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		edgeID := requestPayload.EdgeID
		room := player.Room
		game := room.Game

		err = game.BuildRoad(player.Username, edgeID)
		if err != nil {
			wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
			return true, wsErr
		}

		logs := []string{fmt.Sprintf("%s has built a new road.", player.Username)}

		if game.RoundType() == core.FirstRound {
			room.ProgressStatus()
			logs = append(logs, "Setup phase is over.", "Match starting. Good luck to everyone!")
			room.EnqueueBulkUpdate(
				match.UpdateMapState,
				match.UpdateCurrentRoundPlayerState,
				match.UpdateVertexState,
				match.UpdateEdgeState,
				match.UpdateDiceState,
				match.UpdateLogs(logs),
			)
		} else {
			room.EnqueueBulkUpdate(
				match.UpdateMapState,
				match.UpdateCurrentRoundPlayerState,
				match.UpdateVertexState,
				match.UpdateLogs(logs),
			)
		}

		return true, nil
	default:
		return false, nil
	}

}
