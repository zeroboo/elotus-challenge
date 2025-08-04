package elotuschallenge

import (
	"reflect"
	"testing"
)

func TestSumOfDistancesInTree_Example1_Correct(t *testing.T) {
	n := 6
	edges := [][]int{
		{0, 1}, {0, 2}, {2, 3}, {2, 4}, {2, 5},
	}
	expected := []int{8, 12, 6, 10, 10, 10}
	result := sumOfDistancesInTree(n, edges)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestSumOfDistancesInTree_Example2_Correct(t *testing.T) {
	n := 1
	edges := [][]int{}
	expected := []int{0}
	result := sumOfDistancesInTree(n, edges)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestSumOfDistancesInTree_Example3_Correct(t *testing.T) {
	n := 2
	edges := [][]int{{1, 0}}
	expected := []int{1, 1}
	result := sumOfDistancesInTree(n, edges)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestSumOfDistancesInTree_CompleteGraph_Correct(t *testing.T) {
	n := 4
	edges := [][]int{}
	// Create edges for a complete graph (every node connected to every other node)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges = append(edges, []int{i, j})
		}
	}
	t.Logf("Edges: %v", edges)
	expected := []int{}
	for i := 0; i < n; i++ {
		expected = append(expected, n-1)
	}
	result := sumOfDistancesInTree(n, edges)
	t.Logf("Result: %v", result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}
