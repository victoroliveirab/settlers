package state

import (
	"fmt"

	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) BuildRoad(playerID string, edgeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build road during other player's turn")
		return err
	}

	if state.roundType != SetupRoad1 && state.roundType != SetupRoad2 && state.roundType != Regular {
		err := fmt.Errorf("Cannot build road during %s", RoundTypeTranslation[state.roundType])
		return err
	}

	edge, exists := state.roadMap[edgeID]
	if exists {
		owner := state.findPlayer(edge.Owner)
		err := fmt.Errorf("Player %s already has road at edge #%d", owner, edgeID)
		return err
	}

	if state.roundType == SetupRoad1 || state.roundType == SetupRoad2 {
		if !state.isVertexAllowedSetupPhase(playerID, edgeID) {
			err := fmt.Errorf("Cannot build road in this spot (edge#%d) during setup", edgeID)
			return err
		}
		state.handleNewRoad(playerID, edgeID)
		state.handleChangeSetupRoundType()
		return nil
	}

	resources := state.playerResourceHandMap[playerID]
	if resources["Lumber"] < 1 || resources["Brick"] < 1 {
		err := fmt.Errorf("Insufficient resources to build a road")
		return err
	}

	numberOfRoads := len(state.playerRoadMap[playerID])
	if numberOfRoads >= state.maxRoads {
		err := fmt.Errorf("Cannot have more than %d roads at once", state.maxRoads)
		return err
	}

	if !state.ownsBuildingApproaching(playerID, edgeID) {
		err := fmt.Errorf("Cannot build isolated road (edge#%d)", edgeID)
		return err
	}

	state.playerResourceHandMap[playerID]["Lumber"]--
	state.playerResourceHandMap[playerID]["Brick"]--
	state.handleNewRoad(playerID, edgeID)

	return nil
}

func (state *GameState) isVertexAllowedSetupPhase(playerID string, edgeID int) bool {
	vertexID := utils.SliceLast(state.playerSettlementMap[playerID])
	allowedEdgesIDs := state.definition.EdgesByVertex[vertexID]
	return utils.SliceContains(allowedEdgesIDs, edgeID)
}

func (state *GameState) handleNewRoad(playerID string, edgeID int) {
	entry := Building{
		ID:    edgeID,
		Owner: playerID,
	}
	state.roadMap[edgeID] = entry
	state.playerRoadMap[playerID] = append(state.playerRoadMap[playerID], edgeID)
	state.computeLongestRoad(playerID)
	changed := state.recountLongestRoad()
	if changed {
		state.updatePoints()
	}
}

func (state *GameState) computeLongestRoad(playerID string) {
	graph := make(map[int][]int)
	for _, edgeID := range state.playerRoadMap[playerID] {
		edge := state.definition.VerticesByEdge[edgeID]
		vertex1 := edge[0]
		vertex2 := edge[1]

		_, exists := graph[vertex1]
		if !exists {
			graph[vertex1] = make([]int, 0)
		}
		graph[vertex1] = append(graph[vertex1], edgeID)

		_, exists = graph[vertex2]
		if !exists {
			graph[vertex2] = make([]int, 0)
		}
		graph[vertex2] = append(graph[vertex2], edgeID)
	}

	var maxPath []int
	var dfs func(node int, visited map[int]bool, path []int)

	dfs = func(node int, visited map[int]bool, path []int) {
		if len(path) > len(maxPath) {
			maxPath = append([]int{}, path...)
		}

		for _, edgeID := range graph[node] {
			if !visited[edgeID] {
				var vertex int
				if state.definition.VerticesByEdge[edgeID][0] == node {
					vertex = state.definition.VerticesByEdge[edgeID][1]
				} else if state.definition.VerticesByEdge[edgeID][1] == node {
					vertex = state.definition.VerticesByEdge[edgeID][0]
				} else {
					panic(fmt.Sprintf("unknown edgeID %d", edgeID))
				}
				settlement, settlementExists := state.settlementMap[vertex]
				city, cityExists := state.cityMap[vertex]
				if (settlementExists && settlement.Owner != playerID) || (cityExists && city.Owner != playerID) {
					continue
				}
				visited[edgeID] = true
				dfs(vertex, visited, append(path, edgeID))
				delete(visited, edgeID)
			}
		}
	}

	for startNode := range graph {
		visited := make(map[int]bool)
		dfs(startNode, visited, []int{})
	}

	state.playerLongestRoad[playerID] = maxPath
}
