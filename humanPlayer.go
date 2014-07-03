package main

import "fmt"

type humanPlayer struct {
	char byte
}

func (h *humanPlayer) Input() (int, int, error) {
	var a, b int
	if n, err := fmt.Scanln(&a, &b); err != nil {
		return -1, -1, err
	} else if n == 0 {
		return -1, -1, errBadInput{}
	}
	return a, b, nil
}

func (h *humanPlayer) Char() byte {
	return h.char
}
