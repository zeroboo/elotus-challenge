package graph

import (
	"bytes"
	"fmt"
	"sort"
)

type AdjacentListGraph struct {
	vertices  map[int][]*Edge
	hasWeight bool
}

func NewUnweightedGraph() *AdjacentListGraph {
	graph := &AdjacentListGraph{
		vertices:  make(map[int][]*Edge),
		hasWeight: false,
	}
	return graph
}

func (g *AdjacentListGraph) String() string {
	buf := bytes.Buffer{}

	//Get list of ids by order
	nodeIds := []int{}
	for k, _ := range g.vertices {
		nodeIds = append(nodeIds, k)
	}

	sort.Ints(nodeIds)
	for _, id := range nodeIds {
		vertexEdges := g.vertices[id]
		buf.WriteString(fmt.Sprintf("%d: ", id))
		buf.WriteString("[")
		for _, edge := range vertexEdges {
			if g.hasWeight {
				buf.WriteString(fmt.Sprintf("%d:%v ", edge.Id, edge.Weight))
			} else {
				buf.WriteString(fmt.Sprintf("%d ", edge.Id))
			}
		}
		buf.WriteString("]")
	}

	return buf.String()
}

// Add connection from vertexFrom to vertexTo, if vertexFrom not exists, create new
func (g *AdjacentListGraph) AddWeightedEdge(from int, to int, weight int) {
	vertexFrom := g.vertices[from]
	if vertexFrom == nil {
		vertexFrom = make([]*Edge, 0)
	}
	if !IsEdgeExists(vertexFrom, to) {
		vertexFrom = append(vertexFrom, NewEdge(to, weight))
	}
	g.vertices[from] = vertexFrom
}

// Add connection from vertexFrom to vertexTo, if vertexFrom not exists, create new
func (g *AdjacentListGraph) AddEdge(from int, to int) {
	g.AddWeightedEdge(from, to, 0)
}

// Add new node if not exists
func (g *AdjacentListGraph) AddNewNode(id int) {
	if _, exists := g.vertices[id]; !exists {
		g.vertices[id] = make([]*Edge, 0)
	}
}

func IsEdgeExists(vertices []*Edge, newEdgeId int) bool {
	for _, v := range vertices {
		if v.Id == newEdgeId {
			return true
		}
	}
	return false
}

// Check if childId is a child of parentId
func (g *AdjacentListGraph) IsChildrenNode(childId int, parentId int) bool {
	parentNode := g.vertices[parentId]
	if parentNode == nil {
		return false
	}

	for _, edge := range parentNode {
		if edge.Id == childId {
			return true
		}
	}
	return false
}

func (g *AdjacentListGraph) GetVertexEdges(id int) []*Edge {
	return g.vertices[id]
}

func (g *AdjacentListGraph) GetAllVertexes() map[int][]*Edge {
	return g.vertices
}
