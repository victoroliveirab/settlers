package match

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
)

func handleEndRound(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game
	err := game.EndRound(player.Username)
	if err != nil {
		wsErr := player.WriteJsonError(message.Type, err)
		return true, wsErr
	}

	handleEndRoundResponse(room, player.Username)
	return true, err
}

func handleEndRoundResponse(room *entities.Room, player string) {
	room.EndRound()
	room.StartRound()
	room.StartSubRound(round.BetweenTurns)
	room.EnqueueBulkUpdate(
		UpdateCurrentRoundPlayerState,
		UpdateDiceState,
		UpdatePass,
		UpdateTrade,
		UpdateBuyDevelopmentCard,
		UpdatePlayerDevHandPermissions,
		UpdateLogs([]string{fmt.Sprintf("%s finished their round.", player)}),
	)
}
