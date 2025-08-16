package graph

import "fmt"

type Edge struct {
	Id     int
	Weight int
}

func NewEdge(id int, weight int) *Edge {
	return &Edge{
		Id:     id,
		Weight: weight,
	}
}

func (e *Edge) String() string {
	return fmt.Sprintf("%d:%d", e.Id, e.Weight)
}
