package main

import "math"

type mctsTree struct {
	root *mctsNode
}

type mctsNode struct {
	wins, simulations int
	player            byte
	position          int
	children          []*mctsNode
}

type mctsPlayer struct {
	char byte
	b    *board
	r    *mctsTree
}

// select returns the node to be expanded.
func (m *mctsTree) selectNode() *mctsNode {
	if len(m.root.children) == 0 {
		return m.root
	}
	var n *mctsNode
	var max float64 = 0.0
	for _, v := range m.root.children {
		temp := (float64(v.wins) / float64(v.simulations)) + math.Sqrt2*math.Sqrt(math.Log(float64(m.root.simulations))/float64(v.simulations))
		if temp > max {
			max = temp
			n = v
		}
	}
	return n
}

func (m *mctsPlayer) Char() byte {
	return m.char
}

func (m *mctsPlayer) Input() (int, int, error) {
	//a := getAvailableSpots(m.b)
	return 0, 0, nil
}
