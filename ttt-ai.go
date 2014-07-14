package main

import "fmt"
import "flag"

type player interface {
	Input() (int, int, error)
	Char() byte
}

type errBadInput struct{}

func (e errBadInput) Error() string {
	return "Bad input."
}

type board struct {
	spots [9]byte
}

func (b *board) play(p byte, r, c int) error {
	spot := (r * 3) + c
	if spot < 0 || spot >= len(b.spots) || r >= 3 || c >= 3 {
		return errOffBoard{r, c, "Not a board spot."}
	}
	if b.spots[spot] == 'X' || b.spots[spot] == 'O' {
		return errSpotTaken{r, c, "Spot already taken"}
	}
	b.spots[spot] = p
	return nil
}

func (b *board) checkDone() (bool, byte) {
	for i := 0; i < 3; i++ {
		if b.spots[i] == b.spots[i+3] && b.spots[i] == b.spots[i+6] && b.spots[i] != 0 {
			return true, b.spots[i]
		} else if b.spots[i*3] == b.spots[i*3+1] && b.spots[i*3] == b.spots[i*3+2] && b.spots[i*3] != 0 {
			return true, b.spots[i*3]
		}
	}
	if b.spots[0] == b.spots[4] && b.spots[0] == b.spots[8] && b.spots[0] != 0 {
		return true, b.spots[0]
	} else if b.spots[2] == b.spots[4] && b.spots[2] == b.spots[6] && b.spots[2] != 0 {
		return true, b.spots[2]
	}
	for i := 0; i < len(b.spots); i++ {
		if b.spots[i] == 0 {
			return false, 0
		}
	}
	return true, 0
}

type game struct {
	b      *board
	p1, p2 player
	isDone bool
	turn   int
}

func (g *game) mainLoop() int {
	// g.draw()
	winner := -1
	for !g.isDone {
		if err := g.tick(); err != nil {
			panic(err)
		}
		// g.draw()
		isDone, char := g.b.checkDone()
		if isDone {
			switch char {
			case g.p1.Char():
				winner = 1
			case g.p2.Char():
				winner = 2
			default:
				winner = 0
			}
			g.isDone = isDone
		}
	}
	return winner
}

func (g *game) tick() error {
	switch g.turn {
	case 1:
		doInput(g.p1, g.b)
		g.turn = 2
	case 2:
		doInput(g.p2, g.b)
		g.turn = 1
	default:
		return errBadTurn(g.turn)
	}
	return nil
}

// func (g *game) draw() {
// 	fmt.Println(g.b.spots[0:3])
// 	fmt.Println(g.b.spots[3:6])
// 	fmt.Println(g.b.spots[6:9])
// 	fmt.Println("Player ", g.turn, "'s turn.")
// }

type errBadTurn int

func (e errBadTurn) Error() string {
	return fmt.Sprintf("%v: %v", e, "Not a valid turn value")
}

type errSpotTaken struct {
	r, c int
	what string
}

func (e errSpotTaken) Error() string {
	return fmt.Sprintf("%v, %v: %v", e.r, e.c, e.what)
}

type errOffBoard struct {
	r, c int
	what string
}

func (e errOffBoard) Error() string {
	return fmt.Sprintf("%v, %v: %v", e.r, e.c, e.what)
}

func doInput(p player, b *board) {
	for endTurn := false; !endTurn; {
		r, c, err := p.Input()
		if err == nil {
			if err := b.play(p.Char(), r, c); err != nil {
				fmt.Println(err)
				continue
			}
			endTurn = true
		} else {
			fmt.Println(err)
		}
	}
}

var numGames int

func init() {
	const (
		gamesHelp = "number of games to simulate"
	)
	flag.IntVar(&numGames, "games", 1000, gamesHelp)
	flag.IntVar(&numGames, "g", 1000, gamesHelp+" - shorthand")
}

func main() {
	flag.Parse()
	c := make(chan int)

	for i := 0; i < numGames; i++ {
		go playGame(c)
	}
	s := struct{ tie, p1, p2 int }{0, 0, 0}
	for i := 0; i < numGames; i++ {
		switch <-c {
		case 0:
			s.tie++
		case 1:
			s.p1++
		case 2:
			s.p2++
		}
	}
	fmt.Println(s.tie, s.p1, s.p2)
}

func playGame(c chan int) {
	g := &game{
		b:      &board{},
		isDone: false,
		turn:   1,
	}
	g.p1 = newRandomPlayer('X', g.b)
	g.p2 = newRandomPlayer('O', g.b)
	c <- g.mainLoop()
}
