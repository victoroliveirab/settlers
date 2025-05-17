package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type tileClickRequestPayload struct {
	TileID int `json:"tile"`
}

// REFACTOR: this assumes tile click's only purpose is to move the knight
// If this ever changes, I have to rethink this
func handleTileClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[tileClickRequestPayload](message)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	tileID := payload.TileID
	room := player.Room
	game := room.Game

	err = game.MoveRobber(player.Username, tileID)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	handleMoveRobberResponse(room)
	return true, nil
}

func handleMoveRobberResponse(room *entities.Room) {
	game := room.Game
	currentRoundPlayer := game.CurrentRoundPlayer().ID
	roundType := game.RoundType()
	robbablePlayers, _ := game.RobbablePlayers(currentRoundPlayer)

	logs := make([]string, 0)
	logs = append(logs, fmt.Sprintf("%s moved the robber.", currentRoundPlayer))
	if roundType == round.Regular || roundType == round.BetweenTurns {
		logs = append(logs, fmt.Sprintf("No player to rob."))
		if roundType == round.Regular {
			room.ResumeRound()
		} else {
			// NOTE: this resets the counter instead of keeping in running, but I guess it's fine
			room.StartSubRound(round.BetweenTurns)
		}
	} else if len(robbablePlayers) == 1 {
		game.RobPlayer(currentRoundPlayer, robbablePlayers[0])
		logs = append(logs, fmt.Sprintf("%s robbed %s.", currentRoundPlayer, robbablePlayers[0]))
		if roundType == round.Regular {
			room.ResumeRound()
		} else {
			// NOTE: this resets the counter instead of keeping in running, but I guess it's fine
			room.StartSubRound(round.BetweenTurns)
		}
	} else {
		logs = append(logs, fmt.Sprintf("%s choosing who to rob.", currentRoundPlayer))
		room.StartSubRound(roundType)
	}

	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateMapState,
		UpdatePlayerHand,
		UpdateResourceCount,
		UpdateRobberMovement,
		UpdateRobbablePlayers,
		UpdatePlayerDevHandPermissions,
		UpdatePass,
		UpdateTrade,
		UpdateEdgeState,
		UpdateVertexState,
		UpdateLogs(logs),
	)
}
