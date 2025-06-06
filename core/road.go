package core

import (
	"fmt"

	"github.com/victoroliveirab/settlers/core/packages/round"
	"github.com/victoroliveirab/settlers/utils"
)

func (state *GameState) BuildRoad(playerID string, edgeID int) error {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot build road during other player's turn")
		return err
	}

	// REFACTOR: DONT LIKE THIS HERE
	roundType := state.round.GetRoundType()
	if roundType == round.BuildRoad1Development || roundType == round.BuildRoad2Development {
		return state.PickRoadBuildingSpot(playerID, edgeID)
	}

	if roundType != round.SetupRoad1 && roundType != round.SetupRoad2 && roundType != round.Regular {
		err := fmt.Errorf("Cannot build road during %s", state.round.GetCurrentRoundTypeDescription())
		return err
	}

	edge, exists := state.board.GetRoads()[edgeID]
	if exists {
		owner := state.findPlayer(edge.Owner)
		err := fmt.Errorf("Player %s already has road at edge #%d", owner, edgeID)
		return err
	}

	if roundType == round.SetupRoad1 || roundType == round.SetupRoad2 {
		if !state.isEdgeAllowedSetupPhase(playerID, edgeID) {
			err := fmt.Errorf("Cannot build road in this spot (edge#%d) during setup", edgeID)
			return err
		}
		state.handleNewRoad(playerID, edgeID)
		state.handleChangeSetupRoundType()
		return nil
	}

	playerState := state.playersStates[playerID]

	resources := playerState.GetResources()
	if resources["Lumber"] < 1 || resources["Brick"] < 1 {
		err := fmt.Errorf("Insufficient resources to build a road")
		return err
	}

	numberOfRoads := len(playerState.GetRoads())
	if numberOfRoads >= state.maxRoads {
		err := fmt.Errorf("Cannot have more than %d roads at once", state.maxRoads)
		return err
	}

	if !state.ownsBuildingApproaching(playerID, edgeID) {
		err := fmt.Errorf("Cannot build isolated road (edge#%d)", edgeID)
		return err
	}

	playerState.RemoveResource("Lumber", 1)
	playerState.RemoveResource("Brick", 1)
	state.handleNewRoad(playerID, edgeID)

	return nil
}

func (state *GameState) isEdgeAllowedSetupPhase(playerID string, edgeID int) bool {
	playerState := state.playersStates[playerID]
	vertexID := utils.SliceLast(playerState.GetSettlements())
	allowedEdgesIDs := state.board.Definition.EdgesByVertex[vertexID]
	return utils.SliceContains(allowedEdgesIDs, edgeID)
}

func (state *GameState) handleNewRoad(playerID string, edgeID int) {
	state.board.AddRoad(playerID, edgeID)
	playerState := state.playersStates[playerID]
	playerState.AddRoad(edgeID)

	state.computeLongestRoad(playerID)
	changed := state.recountLongestRoad()
	if changed {
		state.updatePoints()
	}
}

func (state *GameState) computeLongestRoad(playerID string) {
	graph := make(map[int][]int)
	playerState := state.playersStates[playerID]
	for _, edgeID := range playerState.GetRoads() {
		edge := state.board.Definition.VerticesByEdge[edgeID]
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
	settlements := state.board.GetSettlements()
	cities := state.board.GetCities()

	dfs = func(node int, visited map[int]bool, path []int) {
		if len(path) > len(maxPath) {
			maxPath = append([]int{}, path...)
		}

		for _, edgeID := range graph[node] {
			if !visited[edgeID] {
				var vertex int
				if state.board.Definition.VerticesByEdge[edgeID][0] == node {
					vertex = state.board.Definition.VerticesByEdge[edgeID][1]
				} else if state.board.Definition.VerticesByEdge[edgeID][1] == node {
					vertex = state.board.Definition.VerticesByEdge[edgeID][0]
				} else {
					panic(fmt.Sprintf("unknown edgeID %d", edgeID))
				}
				settlement, settlementExists := settlements[vertex]
				city, cityExists := cities[vertex]
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

	playerState.SetLongestRoadSegments(maxPath)
}

func (state *GameState) AvailableEdges(playerID string) ([]int, error) {
	if playerID != state.currentPlayer().ID {
		err := fmt.Errorf("Cannot check available edges during other player's turn")
		return []int{}, err
	}

	roundType := state.round.GetRoundType()
	if roundType != round.SetupRoad1 && roundType != round.SetupRoad2 && roundType != round.Regular && roundType != round.BuildRoad1Development && roundType != round.BuildRoad2Development {
		err := fmt.Errorf("Cannot check available edges during %s", state.round.GetCurrentRoundTypeDescription())
		return []int{}, err
	}

	playerState := state.playersStates[playerID]

	if roundType == round.SetupRoad1 || roundType == round.SetupRoad2 {
		vertexID := utils.SliceLast(playerState.GetSettlements())
		allowedEdgesIDs := state.board.Definition.EdgesByVertex[vertexID]
		return allowedEdgesIDs, nil
	}

	roads := state.board.GetRoads()
	settlements := state.board.GetSettlements()
	edges := utils.NewSet[int]()

	for _, edgeID := range playerState.GetRoads() {
		edge := state.board.Definition.VerticesByEdge[edgeID]
		vertex1 := edge[0]
		vertex2 := edge[1]

		settlement, settlementExistsInVertex1 := settlements[vertex1]
		if !settlementExistsInVertex1 || settlement.Owner == playerID {
			for _, candidateEdgeID := range state.board.Definition.EdgesByVertex[vertex1] {
				_, exists := roads[candidateEdgeID]
				if !exists {
					edges.Add(candidateEdgeID)
				}
			}
		}

		settlement, settlementExistsInVertex2 := settlements[vertex2]
		if !settlementExistsInVertex2 || settlement.Owner == playerID {
			for _, candidateEdgeID := range state.board.Definition.EdgesByVertex[vertex2] {
				_, exists := roads[candidateEdgeID]
				if !exists {
					edges.Add(candidateEdgeID)
				}
			}
		}
	}

	for _, vertexID := range playerState.GetSettlements() {
		for _, candidateEdgeID := range state.board.Definition.EdgesByVertex[vertexID] {
			_, exists := roads[candidateEdgeID]
			if !exists {
				edges.Add(candidateEdgeID)
			}
		}
	}

	for _, vertexID := range playerState.GetCities() {
		for _, candidateEdgeID := range state.board.Definition.EdgesByVertex[vertexID] {
			_, exists := roads[candidateEdgeID]
			if !exists {
				edges.Add(candidateEdgeID)
			}
		}
	}

	return edges.Values(), nil
}
