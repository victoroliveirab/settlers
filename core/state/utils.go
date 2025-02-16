package state

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
