package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func handleEndRound(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game
	err := game.EndRound(player.Username)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateDiceState,
		UpdatePass,
		UpdateTrade,
		UpdateBuyDevelopmentCard,
		UpdatePlayerDevHandPermissions,
		UpdateLogs([]string{fmt.Sprintf("%s finished their round.", player.Username)}),
	)
	return true, err
}
