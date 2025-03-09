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

func TryHandle(player *entities.GamePlayer, message *types.WebSocketMessage) (bool, error) {
	switch message.Type {
	case "game.dice-roll":
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
			// NOTE: should this be a broadcast? Let's rethink this
			room.EnqueueBroadcastMessage(buildMoveRobberDueTo7Broadcast(room), []int64{}, nil)
		} else if game.RoundType() == core.DiscardPhase {
			room.EnqueueBroadcastMessage(buildDiscardCardsBroadcast(room), []int64{}, nil)
		}

		return true, nil
	case "game.discard-cards":
		payload, err := parseDiscardCardsPayload(message.Payload)
		if err != nil {
			wsErr := sendDiscardCardsError(player.Connection, player.ID, err)
			return true, wsErr
		}

		room := player.Room
		game := room.Game
		err = game.DiscardPlayerCards(player.Username, payload.resources)
		if err != nil {
			wsErr := sendDiscardCardsError(player.Connection, player.ID, err)
			return true, wsErr
		}

		formattedResources := utils.FormatResources(payload.resources)
		logs := []string{fmt.Sprintf("%s discarded %s", player.Username, formattedResources)}

		err = sendDiscardCardsSuccess(room, player)
		if err != nil {
			return true, err
		}

		room.EnqueueBroadcastMessage(buildDiscardedCardsBroadcast(room, logs), []int64{}, func() {
			if game.RoundType() == core.MoveRobberDue7 {
				room.EnqueueBroadcastMessage(buildMoveRobberDueTo7Broadcast(room), []int64{}, nil)
			}
		})
		return true, nil
	case "game.new-road":
		payload, err := parseRoadBuildPayload(message.Payload)
		if err != nil {
			wsErr := sendRoadBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		edgeID := payload.edgeID
		room := player.Room
		game := room.Game
		err = game.BuildRoad(player.Username, edgeID)
		if err != nil {
			wsErr := sendRoadBuildError(player.Connection, player.ID, err)
			return true, wsErr
		}

		logs := []string{fmt.Sprintf("%s just built a road.", player.Username)}
		err = sendRoadBuildSuccess(player, edgeID, logs)
		room.EnqueueBroadcastMessage(buildRoadBuildSuccessBroadcast(player.Username, edgeID, logs), []int64{player.ID}, nil)
		return true, err
	case "game.end-round":
		room := player.Room
		game := room.Game
		err := game.EndRound(player.Username)
		if err != nil {
			wsErr := sendEndRoundError(player.Connection, player.ID, err)
			if wsErr != nil {
				return true, wsErr
			}
			return true, nil
		}

		nextPlayer := room.Participants[game.CurrentRoundPlayerIndex()].Player
		err = SendPlayerRound(room, nextPlayer)
		room.EnqueueBroadcastMessage(BuildPlayerRoundOpponentsBroadcast(room), []int64{nextPlayer.ID}, nil)
		return true, err
	default:
		return false, nil
	}
}
