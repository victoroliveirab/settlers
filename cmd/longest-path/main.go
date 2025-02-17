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
		graph[vertex1] = append(graph[vertex1], vertex2)

		_, exists = graph[vertex2]
		if !exists {
			graph[vertex2] = make([]int, 0)
		}
		graph[vertex2] = append(graph[vertex2], vertex1)
	}

	return graph
}

func longestPath(graph map[int][]int) []int {
	var maxPath []int

	var dfs func(node int, visited map[int]bool, path []int)

	dfs = func(node int, visited map[int]bool, path []int) {
		if len(path) > len(maxPath) {
			maxPath = append([]int{}, path...)
		}

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				dfs(neighbor, visited, append(path, neighbor))
				delete(visited, neighbor)
			}
		}
	}

	for startNode := range graph {
		visited := make(map[int]bool)
		visited[startNode] = true
		dfs(startNode, visited, []int{startNode})
	}

	return maxPath

}

func main() {
	s := rand.NewSource(42)
	r := rand.New(s)
	mapsdefinitions.LoadMap("base4")
	def, _ := mapsdefinitions.GenerateMap("base4", r)
	playerEdges := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	graph := makeGraph(def.Definition, playerEdges)

	fmt.Println(graph)
	fmt.Println(longestPath(graph))
}
