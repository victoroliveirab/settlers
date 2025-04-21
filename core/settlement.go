package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/utils"
)

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

	playerState := state.playersStates[playerID]
	resources := playerState.Resources
	if resources["Lumber"] < 1 || resources["Brick"] < 1 || resources["Grain"] < 1 || resources["Sheep"] < 1 {
		err := fmt.Errorf("Insufficient resources to build a settlement")
		return err
	}

	numberOfSettlements := len(playerState.Settlements)
	if numberOfSettlements >= state.maxSettlements {
		err := fmt.Errorf("Cannot have more than %d settlements at once", state.maxSettlements)
		return err
	}

	playerState.RemoveResource("Lumber", 1)
	playerState.RemoveResource("Brick", 1)
	playerState.RemoveResource("Sheep", 1)
	playerState.RemoveResource("Grain", 1)
	state.handleNewSettlement(playerID, vertexID)

	return nil
}

func (state *GameState) handleNewSettlement(playerID string, vertexID int) {
	entry := Building{
		ID:    vertexID,
		Owner: playerID,
	}
	state.settlementMap[vertexID] = entry
	state.playersStates[playerID].AddSettlement(vertexID)

	_, isPort := state.board.Ports[vertexID]
	if isPort {
		state.playersStates[playerID].Ports = append(state.playersStates[playerID].Ports, vertexID)
	}

	// Building a settlement may halt a path
	// OPTIMIZE: check adjacent roads to vertexID and only recalculate for affected players
	for _, player := range state.players {
		state.computeLongestRoad(player.ID)
	}
	state.recountLongestRoad()
	state.updatePoints()
}

func (state *GameState) AvailableVertices(playerID string) ([]int, error) {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot check available vertices during other player's turn")
		return []int{}, err
	}

	if state.roundType != SetupSettlement1 && state.roundType != SetupSettlement2 && state.roundType != Regular {
		err := fmt.Errorf("Cannot check available vertices during %s", RoundTypeTranslation[state.roundType])
		return []int{}, err
	}

	if state.roundType == SetupSettlement1 || state.roundType == SetupSettlement2 {
		availableVertices := make([]int, 0)
		for vertexID := range state.board.Definition.TilesByVertex {
			_, existsSettlement := state.settlementMap[vertexID]
			_, existsCity := state.cityMap[vertexID]
			if existsSettlement || existsCity {
				continue
			}

			blocked := state.isVertexBlocked(vertexID)
			if blocked {
				continue
			}

			availableVertices = append(availableVertices, vertexID)
		}
		return availableVertices, nil
	}

	vertexSet := utils.NewSet[int]()
	for _, edgeID := range state.playersStates[playerID].Roads {
		for _, vertexID := range state.board.Definition.VerticesByEdge[edgeID] {
			_, settlementExists := state.settlementMap[vertexID]
			_, cityExists := state.cityMap[vertexID]
			if settlementExists || cityExists {
				continue
			}

			edgeWithABuilding := state.hasBuildingAtSameEdge(vertexID)
			if edgeWithABuilding == 0 {
				vertexSet.Add(vertexID)
			}
		}
	}

	return vertexSet.Values(), nil
}
