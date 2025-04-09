package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

type pickRobbedPlayerRequestPayload struct {
	Player string `json:"player"`
}

func handlePickRobbedPlayer(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	payload, err := utils.ParseJsonPayload[pickRobbedPlayerRequestPayload](message)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	robbedPlayer := payload.Player
	room := player.Room
	game := room.Game
	err = game.RobPlayer(player.Username, robbedPlayer)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	logs := []string{fmt.Sprintf("%s robbed [res][/res] from %s", player.Username, robbedPlayer)}

	room.EnqueueBulkUpdate(
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdatePass,
		UpdateTrade,
		UpdateRobbablePlayers,
		UpdatePlayerDevHandPermissions,
		UpdateLogs(logs),
	)

	return true, nil
}
