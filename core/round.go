package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) handleChangeSetupRoundType() {
	if state.roundType == SetupSettlement1 || state.roundType == SetupSettlement2 {
		if state.roundType == SetupSettlement1 {
			state.roundType = SetupRoad1
		} else {
			state.roundType = SetupRoad2
		}
		return
	}
	if state.roundType == SetupRoad1 {
		state.currentPlayerIndex++
		if state.currentPlayerIndex == len(state.players) {
			state.roundType = SetupSettlement2
			state.currentPlayerIndex--
		} else {
			state.roundType = SetupSettlement1
		}
		return
	}
	if state.roundType == SetupRoad2 {
		state.currentPlayerIndex--
		if state.currentPlayerIndex < 0 {
			state.roundType = FirstRound
			state.currentPlayerIndex = 0
			state.handOffInitialResources()
		} else {
			state.roundType = SetupSettlement2
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

	if state.roundType != FirstRound && state.roundType != BetweenTurns {
		err := fmt.Errorf("Cannot roll dice during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.dice1 > 0 || state.dice2 > 0 {
		err := fmt.Errorf("Cannot roll dice twice in round")
		return err
	}

	state.dice1 = state.rand.Intn(6) + 1
	state.dice2 = state.rand.Intn(6) + 1

	sum := state.dice1 + state.dice2

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
	state.roundType = Regular
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
		state.roundType = DiscardPhase
	} else {
		state.roundType = MoveRobberDue7
	}
}

func (state *GameState) EndRound(playerID string) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot end round during other player's round")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot end round during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	if state.dice1 == 0 && state.dice2 == 0 {
		err := fmt.Errorf("Cannot end round without rolling dice")
		return err
	}

	state.roundNumber += 1
	state.dice1 = 0
	state.dice2 = 0
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
	state.roundType = BetweenTurns

	state.trade.CancelActiveTrades()
	return nil
}
