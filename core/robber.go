package core

import (
	"fmt"
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
	state.roundType = Regular

	if robberID == robbedID {
		err := fmt.Errorf("Cannot rob from yourself")
		return err
	}

	resources := make([]string, 0)

	for resource, quantity := range state.playerResourceHandMap[robbedID] {
		for i := 0; i < quantity; i++ {
			resources = append(resources, resource)
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

func (state *GameState) UnblockedTiles() []int {
	tileIDs := make([]int, 0)
	for _, tile := range state.tiles {
		if !tile.Blocked {
			tileIDs = append(tileIDs, tile.ID)
		}
	}
	return tileIDs
}
