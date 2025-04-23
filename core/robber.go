package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/utils"
)

// TODO: add prevTileID here
func (state *GameState) MoveRobber(playerID string, tileID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot move robber during other player's turn")
		return err
	}

	roundType := state.round.GetRoundType()
	if roundType != round.MoveRobberDue7 && roundType != round.MoveRobberDueKnight {
		err := fmt.Errorf("Cannot move robber during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	for i, tile := range state.board.GetTiles() {
		if tile.Blocked {
			if tile.ID == tileID {
				err := fmt.Errorf("Cannot move robber to already blocked tile - %d", tileID)
				return err
			}
			state.board.UnblockTileByIndex(i)
			break
		}
	}
	for i, tile := range state.board.GetTiles() {
		if tile.ID == tileID {
			state.board.BlockTileByIndex(i)
			state.round.SetRoundType(round.PickRobbed)
			robbablePlayers, _ := state.RobbablePlayers(playerID)
			if len(robbablePlayers) == 0 {
				// Used knight before round started
				dice := state.round.GetDice()
				if dice[0] == 0 && dice[1] == 0 {
					state.round.SetRoundType(round.BetweenTurns)
				} else {
					state.round.SetRoundType(round.Regular)
				}
			}
			return nil
		}
	}
	// Should never be reached
	err := fmt.Errorf("Cannot move robber to non-existent tile - %d", tileID)
	return err
}

// FIXME: this function is insecure since there's no guarantee that it is moving to the tile it just moved the robber
func (state *GameState) RobPlayer(robberID string, robbedID string) error {
	if robberID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot rob during other player's turn")
		return err
	}

	if state.round.GetRoundType() != round.PickRobbed {
		err := fmt.Errorf("Cannot move robber during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	robbablePlayers, _ := state.RobbablePlayers(robberID)
	if !utils.SliceContains(robbablePlayers, robbedID) {
		err := fmt.Errorf("Cannot rob %s: not in the blocked tile", robbedID)
		return err
	}

	if robberID == robbedID {
		err := fmt.Errorf("Cannot rob from yourself")
		return err
	}

	dice := state.round.GetDice()
	if dice[0] == 0 && dice[1] == 0 {
		state.round.SetRoundType(round.BetweenTurns)
	} else {
		state.round.SetRoundType(round.Regular)
	}

	robbedState := state.playersStates[robbedID]

	resources := make([]string, 0)
	for _, resourceName := range ResourcesOrder {
		quantity := robbedState.Resources[resourceName]
		for i := 0; i < quantity; i++ {
			resources = append(resources, resourceName)
		}
	}

	if len(resources) == 0 {
		err := fmt.Errorf("Cannot rob a player that has no cards")
		return err
	}

	robbedResource := resources[state.rand.Intn(len(resources))]

	robbedState.RemoveResource(robbedResource, 1)
	state.playersStates[robberID].AddResource(robbedResource, 1)
	return nil
}

func (state *GameState) BlockedTiles() []int {
	return state.board.GetBlockedTilesIDs()
}

func (state *GameState) UnblockedTiles() []int {
	return state.board.GetUnblockedTilesIDs()
}
