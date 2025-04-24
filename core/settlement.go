package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) BuildSettlement(playerID string, vertexID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build settlement during other player's turn")
		return err
	}

	roundType := state.round.GetRoundType()
	if roundType != round.SetupSettlement1 && roundType != round.SetupSettlement2 && roundType != round.Regular {
		err := fmt.Errorf("Cannot build settlement during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	vertice, exists := state.board.GetSettlements()[vertexID]
	if exists {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has settlement at vertex #%d", owner.ID, vertexID)
		return err
	}

	vertice, exists = state.board.GetCities()[vertexID]
	if exists {
		owner := state.findPlayer(vertice.Owner)
		err := fmt.Errorf("Player %s already has city at vertex #%d", owner.ID, vertexID)
		return err
	}

	if sharedEdgeID := state.hasBuildingAtSameEdge(vertexID); sharedEdgeID > 0 {
		err := fmt.Errorf("Cannot build at edge %d since it already has a building", sharedEdgeID)
		return err
	}

	if state.round.GetRoundType() == round.SetupSettlement1 || state.round.GetRoundType() == round.SetupSettlement2 {
		state.handleNewSettlement(playerID, vertexID)
		state.handleChangeSetupRoundType()
		return nil
	}

	if !state.ownsRoadApproaching(playerID, vertexID) {
		err := fmt.Errorf("Cannot build at vertex %d since it doesn't have a road attached to it", vertexID)
		return err
	}

	playerState := state.playersStates[playerID]
	resources := playerState.GetResources()
	if resources["Lumber"] < 1 || resources["Brick"] < 1 || resources["Grain"] < 1 || resources["Sheep"] < 1 {
		err := fmt.Errorf("Insufficient resources to build a settlement")
		return err
	}

	numberOfSettlements := len(playerState.GetSettlements())
	if numberOfSettlements >= state.maxSettlements {
		err := fmt.Errorf("Cannot have more than %d settlements at once", state.maxSettlements)
		return err
	}

	playerState.RemoveResource("Lumber", 1)
	playerState.RemoveResource("Brick", 1)
	playerState.RemoveResource("Sheep", 1)
	playerState.RemoveResource("Grain", 1)
	state.bookKeeping.AddResourcesUsed(playerID, "Lumber", 1)
	state.bookKeeping.AddResourcesUsed(playerID, "Brick", 1)
	state.bookKeeping.AddResourcesUsed(playerID, "Sheep", 1)
	state.bookKeeping.AddResourcesUsed(playerID, "Grain", 1)
	state.handleNewSettlement(playerID, vertexID)

	return nil
}

func (state *GameState) handleNewSettlement(playerID string, vertexID int) {
	state.board.AddSettlement(playerID, vertexID)
	state.playersStates[playerID].AddSettlement(vertexID)

	port, isPort := state.board.Ports[vertexID]
	if isPort {
		state.playersStates[playerID].AddPort(vertexID, port)
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

	roundType := state.round.GetRoundType()
	if roundType != round.SetupSettlement1 && roundType != round.SetupSettlement2 && roundType != round.Regular {
		err := fmt.Errorf("Cannot check available vertices during %s", state.round.GetCurrentRoundTypeDescription())
		return []int{}, err
	}

	settlements := state.board.GetSettlements()
	cities := state.board.GetCities()
	if roundType == round.SetupSettlement1 || roundType == round.SetupSettlement2 {
		availableVertices := make([]int, 0)
		for vertexID := range state.board.Definition.TilesByVertex {
			_, existsSettlement := settlements[vertexID]
			_, existsCity := cities[vertexID]
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
	for _, edgeID := range state.playersStates[playerID].GetRoads() {
		for _, vertexID := range state.board.Definition.VerticesByEdge[edgeID] {
			_, settlementExists := settlements[vertexID]
			_, cityExists := cities[vertexID]
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
