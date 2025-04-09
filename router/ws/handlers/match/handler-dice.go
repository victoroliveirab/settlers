package match

import (
	"fmt"
	"maps"

	"github.com/victoroliveirab/settlers/core"
	"github.com/victoroliveirab/settlers/logger"
	"github.com/victoroliveirab/settlers/router/ws/entities"
	"github.com/victoroliveirab/settlers/router/ws/types"
	"github.com/victoroliveirab/settlers/router/ws/utils"
)

func handleDiceRoll(player *entities.GamePlayer, message *types.WebSocketClientRequest) (bool, error) {
	room := player.Room
	game := room.Game

	prevResourceHands := map[string]map[string]int{}
	for _, player := range game.Players() {
		prevResourceHands[player.ID] = maps.Clone(game.ResourceHandByPlayer(player.ID))
	}

	err := game.RollDice(player.Username)
	if err != nil {
		wsErr := utils.WriteJsonError(player.Connection, player.ID, message.Type, err)
		return true, wsErr
	}

	dices := game.Dice()
	dice1 := dices[0]
	dice2 := dices[1]

	logs := make([]string, 1)
	logs[0] = fmt.Sprintf("%s rolled [dice]%d[/dice] [dice]%d[/dice]", player.Username, dice1, dice2)

	for _, player := range game.Players() {
		diff, err := diffResourceHands(prevResourceHands[player.ID], game.ResourceHandByPlayer(player.ID))
		if err != nil {
			logger.LogError(-1, "handleDiceRoll.diffResourceHands", -1, err)
			continue
		}
		if hasDiff(diff) {
			logs = append(logs, fmt.Sprintf("%s got %s", player.ID, serializeHandDiff(diff)))
		}
	}

	if game.RoundType() == core.MoveRobberDue7 {
		logs := []string{fmt.Sprintf("%s moving robber", player.Username)}
		room.EnqueueBulkUpdate(
			UpdateDiceState,
			UpdateRobberMovement,
			UpdatePass,
			UpdateTrade,
			UpdateBuyDevelopmentCard,
			UpdateLogs(logs),
		)
	} else if game.RoundType() == core.DiscardPhase {
		logs := []string{fmt.Sprintf("some players have to discard")}
		room.EnqueueBulkUpdate(
			UpdateDiceState,
			UpdateDiscardPhase,
			UpdatePass,
			UpdateTrade,
			UpdateBuyDevelopmentCard,
			UpdateLogs(logs),
		)
	} else {
		room.EnqueueBulkUpdate(
			UpdateDiceState,
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdatePass,
			UpdateTrade,
			UpdateVertexState,
			UpdateEdgeState,
			UpdateBuyDevelopmentCard,
			UpdateLogs(logs),
		)
	}
	return true, nil
}
