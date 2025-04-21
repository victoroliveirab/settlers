package core

import (
	"fmt"
)

func (state *GameState) BuildCity(playerID string, vertexID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build city during other player's turn")
		return err
	}

	if state.roundType != Regular {
		err := fmt.Errorf("Cannot build city during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	playerState := state.playersStates[playerID]

	if !playerState.HasResourcesToBuildCity() {
		err := fmt.Errorf("Insufficient resources to build a city")
		return err
	}

	numberOfCities := len(playerState.Cities)
	if numberOfCities >= state.maxCities {
		err := fmt.Errorf("Cannot have more than %d cities at once", state.maxCities)
		return err
	}

	playerState.AddCity(vertexID)
	playerState.RemoveResource("Grain", 2)
	playerState.RemoveResource("Ore", 3)
	state.updatePoints()

	return nil
}
