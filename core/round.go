package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) handleChangeSetupRoundType() {
	currentRoundType := state.round.GetRoundType()
	if currentRoundType == round.SetupSettlement1 || currentRoundType == round.SetupSettlement2 {
		if currentRoundType == round.SetupSettlement1 {
			state.round.SetRoundType(round.SetupRoad1)
		} else {
			state.round.SetRoundType(round.SetupRoad2)
		}
		return
	}
	if currentRoundType == round.SetupRoad1 {
		state.currentPlayerIndex++
		if state.currentPlayerIndex == len(state.players) {
			state.round.SetRoundType(round.SetupSettlement2)
			state.currentPlayerIndex--
		} else {
			state.round.SetRoundType(round.SetupSettlement1)
		}
		return
	}
	if currentRoundType == round.SetupRoad2 {
		state.currentPlayerIndex--
		if state.currentPlayerIndex < 0 {
			state.round.SetRoundType(round.FirstRound)
			state.currentPlayerIndex = 0
			state.handOffInitialResources()
		} else {
			state.round.SetRoundType(round.SetupSettlement2)
		}
		return
	}
}

func (state *GameState) handOffInitialResources() {
	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		settlementsIDs := playerState.Settlements
		vertexID := settlementsIDs[1]
		tilesIndexes := state.board.Definition.TilesByVertex[vertexID]
		for _, index := range tilesIndexes {
			tile := state.board.Tiles[index]
			if tile.Resource == "Desert" {
				continue
			}
			playerState.AddResource(tile.Resource, 1)
		}
	}
}

func (state *GameState) RollDice(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot roll dice during other player's turn")
		return err
	}

	roundType := state.round.GetRoundType()
	if roundType != round.FirstRound && roundType != round.BetweenTurns {
		err := fmt.Errorf("Cannot roll dice during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	dice := state.round.GetDice()
	if dice[0] > 0 || dice[1] > 0 {
		err := fmt.Errorf("Cannot roll dice twice in round")
		return err
	}

	dice1 := state.rand.Intn(6) + 1
	dice2 := state.rand.Intn(6) + 1
	state.round.SetDice(dice1, dice2)
	sum := dice1 + dice2

	if sum == 7 {
		state.handle7()
		return nil
	}

	for _, tile := range state.board.Tiles {
		if tile.Token != sum || tile.Blocked || tile.Resource == "Desert" {
			continue
		}
		for _, vertice := range tile.Vertices {
			for _, player := range state.players {
				playerState := state.playersStates[player.ID]
				settlements := playerState.Settlements
				if utils.SliceContains(settlements, vertice) {
					playerState.AddResource(tile.Resource, 1)
				}
				cities := playerState.Cities
				if utils.SliceContains(cities, vertice) {
					playerState.AddResource(tile.Resource, 2)
				}
			}
		}
	}
	state.round.SetRoundType(round.Regular)
	return nil
}

func (state *GameState) handle7() {
	shouldMoveToDiscardPhase := false
	for _, player := range state.players {
		toDiscard := state.discardAmountByPlayer(player.ID)
		if toDiscard > 0 {
			shouldMoveToDiscardPhase = true
			state.playersStates[player.ID].DiscardAmount = toDiscard
		}
	}

	if shouldMoveToDiscardPhase {
		state.round.SetRoundType(round.DiscardPhase)
	} else {
		state.round.SetRoundType(round.MoveRobberDue7)
	}
}

func (state *GameState) EndRound(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot end round during other player's round")
		return err
	}

	if state.round.GetRoundType() != round.Regular {
		err := fmt.Errorf("Cannot end round during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	dice := state.round.GetDice()
	if dice[0] == 0 && dice[1] == 0 {
		err := fmt.Errorf("Cannot end round without rolling dice")
		return err
	}

	state.round.IncrementRound()
	state.round.SetDice(0, 0)
	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		playerState.HasDiscardedThisRound = false
		playerState.NumDevCardsPlayedTurn = 0
		playerState.DiscardAmount = 0
	}
	newIndex := state.currentPlayerIndex + 1
	if newIndex >= len(state.players) {
		newIndex = 0
	}
	state.currentPlayerIndex = newIndex
	state.round.SetRoundType(round.BetweenTurns)

	state.trade.CancelActiveTrades()
	return nil
}
