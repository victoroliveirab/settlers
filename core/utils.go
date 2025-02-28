package core

import (
	"log"

	"github.com/victoroliveirab/settlers/utils"
)

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

func (state *GameState) ownsBuildingApproaching(playerID string, edgeID int) bool {
	vertex1 := state.definition.VerticesByEdge[edgeID][0]
	vertex2 := state.definition.VerticesByEdge[edgeID][1]

	log.Printf("edge %d, vertex1 %d, vertex2 %d\n", edgeID, vertex1, vertex2)

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
