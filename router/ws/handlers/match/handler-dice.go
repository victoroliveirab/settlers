package match

import (
	"fmt"
	"maps"

	"github.com/victoroliveirab/settlers/core/packages/round"
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
	handleDiceRollResponse(room, prevResourceHands)
	return true, nil
}

func handleDiceRollResponse(room *entities.Room, prevResourceHands map[string]map[string]int) {
	game := room.Game
	currentRoundPlayer := game.CurrentRoundPlayer().ID
	dices := game.Dice()
	dice1 := dices[0]
	dice2 := dices[1]

	logs := make([]string, 1)
	logs[0] = fmt.Sprintf("%s rolled [dice v=%d][dice v=%d]", currentRoundPlayer, dice1, dice2)

	for _, player := range game.Players() {
		diff, err := diffResourceHands(prevResourceHands[player.ID], game.ResourceHandByPlayer(player.ID))
		if err != nil {
			logger.LogError(-1, "handleDiceRoll.diffResourceHands", -1, err)
			continue
		}
		if hasDiff(diff) {
			logs = append(logs, fmt.Sprintf("%s got %s", player.ID, formatResourceCollection(diff)))
		}
	}

	if game.RoundType() == round.MoveRobberDue7 {
		logs = append(logs, fmt.Sprintf("%s moving robber", currentRoundPlayer))
		room.StartSubRound(round.MoveRobberDue7)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateDiceState,
			UpdateRobberMovement,
			UpdatePass,
			UpdateTrade,
			UpdateVertexState,
			UpdateEdgeState,
			UpdateBuyDevelopmentCard,
			UpdatePlayerDevHandPermissions,
			UpdateLogs(logs),
		)
	} else if game.RoundType() == round.DiscardPhase {
		logs = append(logs, fmt.Sprintf("some players have to discard"))
		room.StartSubRound(round.DiscardPhase)
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateDiceState,
			UpdateDiscardPhase,
			UpdatePass,
			UpdateTrade,
			UpdateVertexState,
			UpdateEdgeState,
			UpdateBuyDevelopmentCard,
			UpdatePlayerDevHandPermissions,
			UpdateLogs(logs),
		)
	} else {
		room.ResumeRound()
		room.EnqueueBulkUpdate(
			UpdateCurrentRoundPlayerState,
			UpdateDiceState,
			UpdatePlayerHand,
			UpdateResourceCount,
			UpdatePass,
			UpdateTrade,
			UpdateVertexState,
			UpdateEdgeState,
			UpdateBuyDevelopmentCard,
			UpdatePlayerDevHandPermissions,
			UpdateLogs(logs),
		)
	}
}
