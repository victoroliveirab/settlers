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
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	robbedPlayer := payload.Player
	room := player.Room
	game := room.Game
	err = game.RobPlayer(player.Username, robbedPlayer)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	handlePickRobbedResponse(room, player.Username, robbedPlayer)
	return true, nil
}

func handlePickRobbedResponse(room *entities.Room, robber, robbed string) {
	logs := []string{fmt.Sprintf("%s robbed [res q=1 v=?] from %s", robber, robbed)}
	room.ResumeRound()
	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateResourceCount,
		UpdatePlayerHand,
		UpdatePass,
		UpdateTrade,
		UpdateRobbablePlayers,
		UpdatePlayerDevHandPermissions,
		UpdateLogs(logs),
	)
}
