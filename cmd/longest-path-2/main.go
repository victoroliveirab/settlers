package main

import (
	"fmt"
	"math/rand"

	mapsdefinitions "github.com/victoroliveirab/settlers/core/maps"
)

func makeGraph(definition mapsdefinitions.MapDefinition, edges []int) map[int][]int {
	graph := make(map[int][]int)
	for _, edgeID := range edges {
		edge := definition.VerticesByEdge[edgeID]
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

	return graph
}

func longestPath(definition mapsdefinitions.MapDefinition, graph map[int][]int) []int {
	var maxPath []int
	var dfs func(node int, visited map[int]bool, path []int)

	dfs = func(node int, visited map[int]bool, path []int) {
		if len(path) > len(maxPath) {
			maxPath = append([]int{}, path...)
		}

		// fmt.Printf("dfs call - node %d - visited %v - path %v")
		for _, edgeID := range graph[node] {
			if !visited[edgeID] {
				visited[edgeID] = true
				var vertex int
				if definition.VerticesByEdge[edgeID][0] == node {
					vertex = definition.VerticesByEdge[edgeID][1]
				} else if definition.VerticesByEdge[edgeID][1] == node {
					vertex = definition.VerticesByEdge[edgeID][0]
				} else {
					panic(fmt.Sprintf("unknown edgeID %d", edgeID))
				}
				dfs(vertex, visited, append(path, edgeID))
				delete(visited, edgeID)
			}
		}
	}

	for startNode := range graph {
		visited := make(map[int]bool)
		dfs(startNode, visited, []int{})
	}

	return maxPath
}

func main() {
	s := rand.NewSource(42)
	r := rand.New(s)
	mapsdefinitions.LoadMap("base4")
	def, _ := mapsdefinitions.GenerateMap("base4", r)
	playerEdges := []int{1, 2, 3, 4, 5, 6, 7, 19, 22, 23, 24}

	graph := makeGraph(def.Definition, playerEdges)
	fmt.Println(graph)

	longest := longestPath(def.Definition, graph)
	fmt.Println(longest)
}
