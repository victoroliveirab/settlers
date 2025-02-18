package state

import "fmt"

func (state *GameState) BuildSettlement(playerID string, vertexID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build settlement during other player's turn")
		return err
	}

	if state.roundType != SetupSettlement1 && state.roundType != SetupSettlement2 && state.roundType != Regular {
		err := fmt.Errorf("Cannot build settlement during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	vertice, exists := state.settlementMap[vertexID]
	if exists {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has settlement at vertex #%d", owner.ID, vertexID)
		return err
	}

	vertice, exists = state.cityMap[vertexID]
	if exists {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has city at vertex #%d", owner.ID, vertexID)
		return err
	}

	if sharedEdgeID := state.hasBuildingAtSameEdge(vertexID); sharedEdgeID > 0 {
		err := fmt.Errorf("Cannot build at edge %d since it already has a building", sharedEdgeID)
		return err
	}

	if state.roundType == SetupSettlement1 || state.roundType == SetupSettlement2 {
		state.handleNewSettlement(playerID, vertexID)
		state.handleChangeSetupRoundType()
		return nil
	}

	if !state.ownsRoadApproaching(playerID, vertexID) {
		err := fmt.Errorf("Cannot build at vertex %d since it doesn't have a road attached to it", vertexID)
		return err
	}

	resources := state.playerResourceHandMap[playerID]
	if resources["Lumber"] < 1 || resources["Brick"] < 1 || resources["Grain"] < 1 || resources["Sheep"] < 1 {
		err := fmt.Errorf("Insufficient resources to build a settlement")
		return err
	}

	numberOfSettlements := len(state.playerSettlementMap[playerID])
	if numberOfSettlements >= state.maxSettlements {
		err := fmt.Errorf("Cannot have more than %d settlements at once", state.maxSettlements)
		return err
	}

	state.playerResourceHandMap[playerID]["Lumber"]--
	state.playerResourceHandMap[playerID]["Brick"]--
	state.playerResourceHandMap[playerID]["Sheep"]--
	state.playerResourceHandMap[playerID]["Grain"]--
	state.handleNewSettlement(playerID, vertexID)

	return nil
}

func (state *GameState) handleNewSettlement(playerID string, vertexID int) {
	entry := Building{
		ID:    vertexID,
		Owner: playerID,
	}
	state.settlementMap[vertexID] = entry
	state.playerSettlementMap[playerID] = append(state.playerSettlementMap[playerID], vertexID)

	// Building a settlement may halt a path
	// OPTIMIZE: check adjacent roads to vertexID and only recalculate for affected players
	for _, player := range state.players {
		state.computeLongestRoad(player.ID)
	}
}
