package main

import "fmt"
import "flag"
import "runtime"

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
		if b.spots[i] == b.spots[i+3] &&
			b.spots[i] == b.spots[i+6] &&
			b.spots[i] != 0 {
			return true, b.spots[i]
		} else if b.spots[i*3] == b.spots[i*3+1] &&
			b.spots[i*3] == b.spots[i*3+2] &&
			b.spots[i*3] != 0 {
			return true, b.spots[i*3]
		}
	}
	if b.spots[0] == b.spots[4] &&
		b.spots[0] == b.spots[8] &&
		b.spots[0] != 0 {
		return true, b.spots[0]
	} else if b.spots[2] == b.spots[4] &&
		b.spots[2] == b.spots[6] &&
		b.spots[2] != 0 {
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
	turn   player
}

func (g *game) mainLoop() byte {
	g.draw()
	var winner byte = 0
	for !g.isDone {
		isDone, char := g.b.checkDone()
		if isDone {
			switch char {
			case g.p1.Char():
				winner = g.p1.Char()
			case g.p2.Char():
				winner = g.p2.Char()
			default:
				winner = 'T'
			}
			g.isDone = isDone
			continue
		}
		if err := g.tick(); err != nil {
			panic(err)
		}
		g.draw()
	}
	return winner
}

func (g *game) tick() error {
	if g.turn != g.p1 && g.turn != g.p2 {
		return errBadTurn(g.turn.Char())
	}
	doInput(g.turn, g.b)
	if g.turn == g.p1 {
		g.turn = g.p2
	} else {
		g.turn = g.p1
	}
	return nil
}

func (g *game) draw() {
	fmt.Println(g.b.spots[0:3])
	fmt.Println(g.b.spots[3:6])
	fmt.Println(g.b.spots[6:9])
	fmt.Println("Player ", g.turn, "'s turn.")
}

type errBadTurn byte

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
	fmt.Println(runtime.GOMAXPROCS(8))
	c := make(chan byte)

	for i := 0; i < numGames; i++ {
		go playGame(c)
	}
	s := struct{ tie, p1, p2 int }{0, 0, 0}
	for i := 0; i < numGames; i++ {
		switch <-c {
		case 'T':
			s.tie++
		case 'X':
			s.p1++
		case 'O':
			s.p2++
		}
	}
	fmt.Println(s.tie, s.p1, s.p2)
}

func playGame(c chan byte) {
	g := &game{
		b:      &board{},
		isDone: false,
	}
	g.p1 = &humanPlayer{'X'}
	g.p2 = newMctsPlayer('O', g.b)
	g.turn = g.p1
	c <- g.mainLoop()
}
