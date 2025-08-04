package elotuschallenge

import (
	"fmt"
	"log"
	"strings"
)

type Node struct {
	Value    int
	Children []*Node
}

func (n *Node) ToString() string {
	if n == nil {
		return "NilNode"
	}
	var buf strings.Builder

	for _, child := range n.Children {
		if child != nil {
			buf.WriteString(fmt.Sprintf("%v ", child.ToString()))
		}
	}

	return fmt.Sprintf("%d[%v]", n.Value, buf.String())
}

func (n *Node) FindChildNode(value int) *Node {
	for _, child := range n.Children {
		if child.Value == value {
			return child
		}
	}
	for _, child := range n.Children {
		return child.FindChildNode(value)
	}
	return nil
}

func NewNode(value int) *Node {
	return &Node{
		Value:    value,
		Children: []*Node{},
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

type IAccumulateFunc interface {
	Calculate(node *Node, depth int) int
	GetValue() int
}

type DepthMeasurer struct {
	TotalDepth int
}

func NewDepthMeasurer(depth int) *DepthMeasurer {
	return &DepthMeasurer{
		TotalDepth: depth,
	}
}

func (cc *DepthMeasurer) Calculate(node *Node, depth int) int {
	cc.TotalDepth += depth
	return cc.TotalDepth
}

func (cc *DepthMeasurer) GetValue() int {
	return cc.TotalDepth
}

type Visitor struct {
	Visisted map[int]map[int]bool
	AccFunc  IAccumulateFunc
	Depth    int
	Root     *Node
}

func NewVisitor(accFunc IAccumulateFunc) *Visitor {
	return &Visitor{
		Visisted: make(map[int]map[int]bool),
		AccFunc:  accFunc,
		Depth:    0,
	}
}
func (v *Visitor) MarkVisited(node *Node, depth int) {
	if v.Visisted[node.Value] == nil {
		v.Visisted[node.Value] = make(map[int]bool)
	}
	v.Visisted[node.Value][depth] = true
}
func (v *Visitor) IsVisited(node *Node, depth int) bool {
	if v.Visisted[node.Value] == nil {
		return false
	}
	return v.Visisted[node.Value][depth]
}
func (v *Visitor) Visit(node *Node, depth int) {
	log.Printf("Begin Visiting node: %v, value=%v, depth=%d", node.ToString(), v.AccFunc.GetValue(), depth)
	v.AccFunc.Calculate(node, depth)
	for _, child := range node.Children {
		if child != v.Root {
			v.Visit(child, depth+1)
		}
	}
	log.Printf("Finish Visiting node: %v, value=%v, depth=%d", node.ToString(), v.AccFunc.GetValue(), depth)
}
func (v *Visitor) Start(node *Node) {

	if node == nil {
		return
	}
	v.Root = node
	v.Visit(node, 0)

}
