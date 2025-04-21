package core

import (
	"math"

	"github.com/victoroliveirab/settlers/utils"
)

var ResourcesOrder [5]string = [5]string{"Lumber", "Brick", "Sheep", "Grain", "Ore"}

func (state *GameState) NumberOfResourcesByPlayer() map[string]int {
	resourcesByPlayer := make(map[string]int)
	for player, pState := range state.playersStates {
		resourcesByPlayer[player] = 0
		for _, count := range pState.Resources {
			resourcesByPlayer[player] += count
		}
	}
	return resourcesByPlayer
}

func (state *GameState) NumberOfDevCardsByPlayer() map[string]int {
	devCardsByPlayer := make(map[string]int)
	for player, pState := range state.playersStates {
		devCardsByPlayer[player] = 0
		for _, cards := range pState.DevelopmentCards {
			devCardsByPlayer[player] += len(cards)
		}
	}
	return devCardsByPlayer
}

func (state *GameState) discardAmountByPlayer(playerID string) int {
	total := 0
	for _, count := range state.playersStates[playerID].Resources {
		total += count
	}
	if total <= state.maxCards {
		return 0
	}
	return int(math.Floor(float64(total) / 2))
}

func (state *GameState) hasBuildingAtSameEdge(vertexID int) int {
	edges := state.board.Definition.EdgesByVertex[vertexID]
	for _, edgeID := range edges {
		vertices := state.board.Definition.VerticesByEdge[edgeID]
		vertex1 := vertices[0]
		vertex2 := vertices[1]

		var vertex int
		if vertex1 == vertexID {
			vertex = vertex2
		} else {
			vertex = vertex1
		}

		_, hasSettlement := state.settlementMap[vertex]
		_, hasCity := state.cityMap[vertex]
		if hasSettlement || hasCity {
			return edgeID
		}
	}
	return 0
}

func (state *GameState) isVertexBlocked(vertexID int) bool {
	edgeID := state.hasBuildingAtSameEdge(vertexID)
	return edgeID > 0
}

func (state *GameState) ownsBuildingApproaching(playerID string, edgeID int) bool {
	vertex1 := state.board.Definition.VerticesByEdge[edgeID][0]
	vertex2 := state.board.Definition.VerticesByEdge[edgeID][1]

	playerState := state.playersStates[playerID]
	hasSettlementVertex1 := utils.SliceContains(playerState.Settlements, vertex1)
	hasSettlementVertex2 := utils.SliceContains(playerState.Settlements, vertex2)

	if hasSettlementVertex1 || hasSettlementVertex2 {
		return true
	}

	hasRoadApproachingVertex1 := state.ownsRoadApproaching(playerID, vertex1)
	if hasRoadApproachingVertex1 {
		return true
	}

	hasRoadApproachingVertex2 := state.ownsRoadApproaching(playerID, vertex2)
	if hasRoadApproachingVertex2 {
		return true
	}

	return false
}

func (state *GameState) ownsRoadApproaching(playerID string, vertexID int) bool {
	edges := state.board.Definition.EdgesByVertex[vertexID]
	for _, edgeID := range edges {
		road, exists := state.roadMap[edgeID]
		if exists && road.Owner == playerID {
			return true
		}
	}
	return false
}
