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
	tiles := state.board.GetTiles()
	for _, player := range state.players {
		playerState := state.playersStates[player.ID]
		settlementsIDs := playerState.GetSettlements()
		vertexID := settlementsIDs[1]
		tilesIndexes := state.board.Definition.TilesByVertex[vertexID]
		for _, index := range tilesIndexes {
			tile := tiles[index]
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
	state.bookKeeping.AddDiceEntry(playerID, sum)

	if sum == 7 {
		state.handle7()
		return nil
	}

	for _, tile := range state.board.GetTiles() {
		if tile.Token != sum || tile.Resource == "Desert" {
			continue
		}
		for _, vertice := range tile.Vertices {
			for _, player := range state.players {
				playerState := state.playersStates[player.ID]
				settlements := playerState.GetSettlements()
				if utils.SliceContains(settlements, vertice) {
					if tile.Blocked {
						state.bookKeeping.AddResourcesBlocked(player.ID, tile.Resource, 1)
					} else {
						playerState.AddResource(tile.Resource, 1)
						state.bookKeeping.AddResourceDrawn(player.ID, tile.Resource, 1)
					}
				}
				cities := playerState.GetCities()
				if utils.SliceContains(cities, vertice) {
					if tile.Blocked {
						state.bookKeeping.AddResourcesBlocked(player.ID, tile.Resource, 2)
					} else {
						playerState.AddResource(tile.Resource, 2)
						state.bookKeeping.AddResourceDrawn(player.ID, tile.Resource, 2)
					}
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
			state.playersStates[player.ID].SetDiscardAmount(toDiscard)
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
		playerState.SetHasDiscardedThisTurn(false)
		playerState.ResetNumberOfDevCardsPlayedCurrentTurn()
		playerState.SetDiscardAmount(0)
	}
	newIndex := state.currentPlayerIndex + 1
	if newIndex >= len(state.players) {
		newIndex = 0
	}
	state.currentPlayerIndex = newIndex
	state.bookKeeping.AddPointsRecord(state.points)
	state.bookKeeping.AddLongestRoadRecord(state.LongestRoadLengths())
	state.round.SetRoundType(round.BetweenTurns)

	state.trade.CancelActiveTrades()
	return nil
}
