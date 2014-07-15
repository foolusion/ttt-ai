package main

import "math/rand"
import "time"

type randomPlayer struct {
	char byte
	b    *board
	rn   *rand.Rand
}

func newRandomPlayer(char byte, b *board) *randomPlayer {
	return &randomPlayer{
		char: char,
		b:    b,
		rn:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *randomPlayer) Input() (int, int, error) {
	a := getAvailableSpots(r.b)
	if len(a) <= 1 {
		return a[0] / 3, a[0] % 3, nil
	}
	i := r.rn.Intn(len(a))
	return a[i] / 3, a[i] % 3, nil
}

func (r *randomPlayer) Char() byte {
	return r.char
}

func getAvailableSpots(b *board) []int {
	available := make([]int, 0, 9)
	for i := range b.spots {
		if b.spots[i] == 0 {
			available = append(available, i)
		}
	}
	return available
}
