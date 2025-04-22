package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
)

func (state *GameState) BuildCity(playerID string, vertexID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build city during other player's turn")
		return err
	}

	if state.round.GetRoundType() != round.Regular {
		err := fmt.Errorf("Cannot build city during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	vertice, exists := state.board.Settlements[vertexID]
	if exists && vertice.Owner != playerID {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has settlement at vertex #%d", owner.ID, vertexID)
		return err
	}

	vertice, exists = state.board.Cities[vertexID]
	if exists {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has city at vertex #%d", owner.ID, vertexID)
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

	state.board.AddCity(playerID, vertexID)
	playerState.AddCity(vertexID)
	playerState.RemoveResource("Grain", 2)
	playerState.RemoveResource("Ore", 3)
	state.updatePoints()

	return nil
}
