package graph

import (
	"bytes"
	"fmt"
)

type Node struct {
	Id       int
	Children []*Node
	Weight   int
}

func NewNode(id int) *Node {
	return &Node{
		Id:       id,
		Children: make([]*Node, 0),
		Weight:   0,
	}
}

func NewNodeWeight(id int, weight int) *Node {
	return &Node{
		Id:       id,
		Children: make([]*Node, 0),
		Weight:   weight,
	}
}

func (n *Node) GetId() int {
	return n.Id
}

func (n *Node) GetWeight() int {
	return n.Weight
}

func (n *Node) HasChild(id int) bool {
	for _, child := range n.Children {
		if child.Id == id {
			return true
		}
	}
	return false

}
func (n *Node) RemoveChild(id int) *Node {
	removeIndex := -1
	for i, child := range n.Children {
		if child.Id == id {
			removeIndex = i
			break
		}
	}
	if removeIndex != -1 {
		n.Children = append(n.Children[:removeIndex], n.Children[removeIndex+1:]...)
	}
	return nil
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) String() string {
	var buffer bytes.Buffer
	if n.Weight != 0 {
		buffer.WriteString(fmt.Sprintf("Node{%d|%d", n.Id, n.Weight))
	} else {
		buffer.WriteString(fmt.Sprintf("Node{%d", n.Id))
	}
	buffer.WriteString("[")
	for i, child := range n.Children {
		buffer.WriteString(fmt.Sprintf("%d", child.Id))
		if i < len(n.Children)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]}")
	return buffer.String()
}
