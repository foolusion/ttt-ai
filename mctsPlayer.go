package main

import (
	"math"
	"math/rand"
	"time"
)

type mctsTree struct {
	root *mctsNode
}

func newMctsTree(player byte) *mctsTree {
	var opponent byte
	if player == 'X' {
		opponent = 'O'
	} else {
		opponent = 'X'
	}
	root := &mctsNode{
		wins:        0,
		simulations: 0,
		player:      opponent,
		parent:      nil,
	}
	return &mctsTree{root}
}

type mctsNode struct {
	wins, simulations int
	player            byte
	position          int
	parent            *mctsNode
	children          []*mctsNode
}

type mctsPlayer struct {
	char byte
	b    *board
	r    *mctsTree
	rn   *rand.Rand
}

func newMctsPlayer(c byte, b *board) *mctsPlayer {
	return &mctsPlayer{c, b, nil, rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// selectNode returns the node to be expanded.
func (m *mctsPlayer) selectNode(n *mctsNode) *mctsNode {
	if n.children == nil {
		return n
	}
	if len(n.children) == 0 {
		return n
	}
	if len(n.children) == 1 {
		return n.children[0]
	}
	var next *mctsNode // next node to select value
	var max float64 = 0.0
	for i, v := range n.children {
		temp := (float64(v.wins) / float64(v.simulations)) +
			math.Sqrt2*math.Sqrt(
				math.Log(float64(n.simulations))/float64(v.simulations))
		if temp > max {
			max = temp
			next = n.children[i]
		}
	}
	if next == nil {
		next = n.children[m.rn.Intn(len(n.children))]
	}
	return m.selectNode(next)
}

func (m *mctsPlayer) expandNode(n *mctsNode) *mctsNode {
	a := getAvailableSpots(m.b)
	var player byte
	if n.player == 'X' {
		player = 'O'
	} else {
		player = 'X'
	}
	n.children = make([]*mctsNode, len(a))
	for i, v := range a {
		n.children[i] = &mctsNode{
			wins:        0,
			simulations: 0,
			player:      player,
			position:    v,
			parent:      n,
		}
	}
	return n.children[m.rn.Intn(len(n.children))]
}

func (m *mctsPlayer) simulateNode(n *mctsNode) bool {
	g := &game{
		b:      &board{},
		isDone: false,
	}
	*g.b = *m.b
	for cur := n; cur.parent != nil; cur = cur.parent {
		g.b.spots[cur.position] = cur.player
	}
	g.p1 = newRandomPlayer('X', g.b)
	g.p2 = newRandomPlayer('O', g.b)
	if n.player == g.p1.Char() {
		g.turn = g.p2
	} else {
		g.turn = g.p1
	}
	winner := g.mainLoop()
	if winner == m.char {
		n.wins++
		n.simulations++
		return true
	}
	n.simulations++
	return false
}

func (m *mctsPlayer) backPropogateNode(n *mctsNode, win bool) {
	for cur := n; cur.parent != nil; cur = cur.parent {
		if win {
			n.parent.wins++
		}
		n.parent.simulations++
	}
}

func (m *mctsPlayer) Char() byte {
	return m.char
}

func (m *mctsPlayer) Input() (int, int, error) {
	m.r = newMctsTree(m.char)
	position := mctsLoop(m)
	return position / 3, position % 3, nil
}

func mctsLoop(m *mctsPlayer) int {
	timer := time.NewTimer(10 * time.Second)
	for i := 0; i < 10000; {
		select {
		case <-timer.C:
			break
		default:
		  mctsTick(m)
		  i++
		}
	}
	max, pos := 0, 0
	for _, v := range m.r.root.children {
		if v.wins >= max {
			max = v.wins
			pos = v.position
		}
	}
	return pos
}

func mctsTick(m *mctsPlayer) {
	n := m.selectNode(m.r.root)
	n = m.expandNode(n)
	win := m.simulateNode(n)
	m.backPropogateNode(n, win)
}
