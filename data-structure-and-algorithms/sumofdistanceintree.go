package elotuschallenge

func sumOfDistancesInTree(n int, edges [][]int) []int {
	// Build adjacency list from edges
	graph := make([][]int, n)
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
	}

	totalDistanceOfNodes := make([]int, n)

	// For each node, calculate sum of distances to all other nodes
	for i := 0; i < n; i++ {
		visited := make([]bool, n)
		totalDistance := 0

		// BFS from node i to calculate distances to all other nodes
		queue := []int{i}
		distances := make([]int, n)
		visited[i] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			for _, neighbor := range graph[current] {
				if !visited[neighbor] {
					visited[neighbor] = true
					distances[neighbor] = distances[current] + 1
					totalDistance += distances[neighbor]
					queue = append(queue, neighbor)
				}
			}
		}

		totalDistanceOfNodes[i] = totalDistance
	}

	return totalDistanceOfNodes
}
