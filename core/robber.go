package core

import (
	"fmt"
	"sort"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) MoveRobber(playerID string, tileID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot move robber during other player's turn")
		return err
	}

	if state.roundType != MoveRobberDue7 && state.roundType != MoveRobberDueKnight {
		err := fmt.Errorf("Cannot move robber during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	for _, tile := range state.tiles {
		if tile.Blocked {
			if tile.ID == tileID {
				err := fmt.Errorf("Cannot move robber to already blocked tile - %d", tileID)
				return err
			}
			tile.Blocked = false
			break
		}
	}
	for _, tile := range state.tiles {
		if tile.ID == tileID {
			tile.Blocked = true
			state.roundType = PickRobbed
			robbablePlayers, _ := state.RobbablePlayers(playerID)
			if len(robbablePlayers) == 0 {
				state.roundType = Regular
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

	if state.roundType != PickRobbed {
		err := fmt.Errorf("Cannot move robber during %s", RoundTypeTranslation[state.roundType])
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

	if state.dice1 == 0 && state.dice2 == 0 {
		state.roundType = BetweenTurns
	} else {
		state.roundType = Regular
	}

	// NOTE: this is done to enforce ordering for tests (math/random seed)
	keys := make([]string, 0, len(state.playerResourceHandMap[robbedID]))
	for kind := range state.playerResourceHandMap[robbedID] {
		keys = append(keys, kind)
	}
	sort.Strings(keys)

	resources := make([]string, 0)
	for _, resourceName := range keys {
		quantity := state.playerResourceHandMap[robbedID][resourceName]
		for i := 0; i < quantity; i++ {
			resources = append(resources, resourceName)
		}
	}

	if len(resources) == 0 {
		err := fmt.Errorf("Cannot rob a player that has no cards")
		return err
	}

	robbedResource := resources[state.rand.Intn(len(resources))]

	state.playerResourceHandMap[robberID][robbedResource]++
	state.playerResourceHandMap[robbedID][robbedResource]--
	return nil
}

func (state *GameState) BlockedTiles() []int {
	tileIDs := make([]int, 0)
	for _, tile := range state.tiles {
		if tile.Blocked {
			tileIDs = append(tileIDs, tile.ID)
		}
	}
	return tileIDs
}

func (state *GameState) UnblockedTiles() []int {
	tileIDs := make([]int, 0)
	for _, tile := range state.tiles {
		if !tile.Blocked {
			tileIDs = append(tileIDs, tile.ID)
		}
	}
	return tileIDs
}
