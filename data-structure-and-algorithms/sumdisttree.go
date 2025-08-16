package elotuschallenge

import (
	"fmt"

	"elotuschallenge/graph"
)

func SumDistanceInTree(totalNodes int, edges [][]int) []int {

	//Build graphs
	trees := []*graph.AdjacentListGraph{}
	for i := 0; i < totalNodes; i++ {
		tree := graph.NewUnweightedGraph()
		trees = append(trees, tree)

		for j := 0; j < totalNodes; j++ {
			tree.AddNewNode(j)
		}
	}

	//Setup nodes for graph

	//Setup edges
	for _, edge := range edges {
		if len(edge) != 2 {
			fmt.Println("Invalid edge length, expected 2, got", len(edge))
			continue
		}
		from := edge[0]
		to := edge[1]
		for _, tree := range trees {
			tree.AddEdge(from, to)
			tree.AddEdge(to, from)
		}
	}

	//Trace all tree
	for i, tree := range trees {
		fmt.Println("Tree", i, ":", tree)

	}

	//Count all connected edges from each node
	result := make([]int, totalNodes)
	for i := 0; i < totalNodes; i++ {
		tree := trees[i]
		result[i] = SumAllDistancesInTree(i, tree)
	}

	return result
}

func SumAllDistancesInTree(root int, tree *graph.AdjacentListGraph) int {
	fmt.Println("--- SumAllDistancesInTree", root)
	fmt.Println("Visiting tree", tree)
	allDistance := 0

	for node := range tree.GetAllVertexes() {
		if node != root {
			//Find path from root to this node
			stack := []int{}
			found := false

			currentNode := root
			visitPath := []int{}
			visited := make(map[int]bool)

			for {
				if len(visitPath) == 0 {
					visitPath = append(visitPath, currentNode)
				} else {
					for {
						previousNodePathId := visitPath[len(visitPath)-1]
						if tree.IsChildrenNode(currentNode, previousNodePathId) {
							visitPath = append(visitPath, currentNode)
							break
						}
						visitPath = visitPath[:len(visitPath)-1]
						if len(visitPath) == 0 {
							visitPath = append(visitPath, currentNode)
							break
						}
					}
				}

				edges := tree.GetVertexEdges(currentNode)
				visited[currentNode] = true
				found = currentNode == node
				if found {
					fmt.Println("Found path from", root, "to", node, ":", visitPath)
					break
				} else {
					for _, edge := range edges {
						if !visited[edge.Id] {
							stack = append(stack, edge.Id)
						}
					}
				}

				if len(stack) == 0 {
					break
				}

				//Pop from stack
				currentNode = stack[len(stack)-1]
				stack = stack[:len(stack)-1]

			}
			distance := 0
			if found {
				distance += len(visitPath) - 1
			}

			fmt.Println("Total distance from", root, "to", node, "is", distance, "path", visitPath)
			allDistance += distance
		}
	}

	return allDistance
}
