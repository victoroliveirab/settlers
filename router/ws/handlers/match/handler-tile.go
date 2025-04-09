package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type tileClickRequestPayload struct {
	TileID int `json:"tile"`
}

func handleTileClick(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[tileClickRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	tileID := payload.TileID
	room := player.Room
	game := room.Game
	logs := make([]string, 0)

	err = game.MoveRobber(player.Username, tileID)
	if err == nil {
		robbablePlayers, _ := game.RobbablePlayers(player.Username)
		if len(robbablePlayers) == 0 {
			logs = append(logs, fmt.Sprintf("No player to rob."))
			room.EnqueueBulkUpdate(
				UpdateMapState,
				UpdatePass,
				UpdateTrade,
				UpdateLogs(logs),
			)
			return true, nil
		} else if len(robbablePlayers) == 1 {
			err := game.RobPlayer(player.Username, robbablePlayers[0])
			if err != nil { // The only player in tile (besides the current player) doesn't have any resources
				// NOTE: This should never happen since the core logic
				logs = append(logs, fmt.Sprintf("No player to rob."))
			} else {
				logs = append(logs, fmt.Sprintf("%s robbed %s.", player.Username, robbablePlayers[0]))
			}
			room.EnqueueBulkUpdate(
				UpdateMapState,
				UpdatePlayerHand,
				UpdateRobberMovement,
				UpdateResourceCount,
				UpdatePass,
				UpdateTrade,
				UpdateLogs(logs),
			)
		} else {
			logs = append(logs, fmt.Sprintf("%s choosing who to rob.", player.Username))
			room.EnqueueBulkUpdate(
				UpdateMapState,
				UpdateRobberMovement,
				UpdateRobbablePlayers,
				UpdatePass,
				UpdateTrade,
				UpdateLogs(logs),
			)
		}
		return true, nil
	}
	// Implemented this way if there are more actions to account for in the tile click

	wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
	return true, wsErr
}
