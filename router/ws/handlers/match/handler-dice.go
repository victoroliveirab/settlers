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

func handleDiceRoll(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	room := player.Room
	game := room.Game

	prevResourceHands := map[string]map[string]int{}
	for _, player := range game.Players() {
		prevResourceHands[player.ID] = maps.Clone(game.ResourceHandByPlayer(player.ID))
	}

	err := game.RollDice(player.Username)
	if err != nil {
		logger.LogError(player.ID, "engine.state.RollDices", -1, err)
		wsErr := sendDiceRollError(player.Connection, player.ID, err)
		if wsErr != nil {
			return true, wsErr
		}
		return true, nil
	}

	dices := game.Dice()
	dice1 := dices[0]
	dice2 := dices[1]

	logs := make([]string, 1)
	logs[0] = fmt.Sprintf("%s rolled [dice]%d[/dice] [dice]%d[/dice]", player.Username, dice1, dice2)

	for _, player := range game.Players() {
		diff, err := diffResourceHands(prevResourceHands[player.ID], game.ResourceHandByPlayer(player.ID))
		if err != nil {
			return true, err
		}
		if hasDiff(diff) {
			logs = append(logs, fmt.Sprintf("%s got %s", player.ID, serializeHandDiff(diff)))
		}
	}

	for _, participant := range room.Participants {
		// TODO: handle error properly
		diceRollMessage := buildDiceRollSuccess(room, participant.Player, logs)
		// TODO: this probably blocks the main thread, so we should spin up go routines later
		utils.WriteJson(participant.Player.Connection, participant.Player.ID, diceRollMessage)
	}

	if game.RoundType() == core.MoveRobberDue7 {
		logs := []string{fmt.Sprintf("%s moving robber", player.Username)}
		// NOTE: should this be a broadcast? Let's rethink this
		room.EnqueueBroadcastMessage(buildMoveRobberBroadcast(room, logs), []int64{}, nil)
	} else if game.RoundType() == core.DiscardPhase {
		room.EnqueueBroadcastMessage(buildDiscardCardsBroadcast(room), []int64{}, nil)
	}

	return true, nil
}
