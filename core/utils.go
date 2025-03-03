package core

import (
	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) NumberOfResourcesByPlayer() map[string]int {
	resourcesByPlayer := make(map[string]int)
	for player, hand := range state.playerResourceHandMap {
		resourcesByPlayer[player] = 0
		for _, count := range hand {
			resourcesByPlayer[player] += count
		}
	}
	return resourcesByPlayer
}

func (state *GameState) NumberOfDevCardsByPlayer() map[string]int {
	devCardsByPlayer := make(map[string]int)
	for player, hand := range state.playerDevelopmentHandMap {
		devCardsByPlayer[player] = 0
		for _, count := range hand {
			devCardsByPlayer[player] += count
		}
	}
	return devCardsByPlayer
}

func (state *GameState) hasBuildingAtSameEdge(vertexID int) int {
	edges := state.definition.EdgesByVertex[vertexID]
	for _, edgeID := range edges {
		vertices := state.definition.VerticesByEdge[edgeID]
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
	vertex1 := state.definition.VerticesByEdge[edgeID][0]
	vertex2 := state.definition.VerticesByEdge[edgeID][1]

	hasSettlementVertex1 := utils.SliceContains(state.playerSettlementMap[playerID], vertex1)
	hasSettlementVertex2 := utils.SliceContains(state.playerSettlementMap[playerID], vertex2)

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
	edges := state.definition.EdgesByVertex[vertexID]
	for _, edgeID := range edges {
		road, exists := state.roadMap[edgeID]
		if exists && road.Owner == playerID {
			return true
		}
	}
	return false
}
