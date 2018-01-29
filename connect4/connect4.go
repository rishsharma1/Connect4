package main

import (
	"errors"
	"fmt"
)

type game struct {
	board   []string
	columns int
	rows    int
}

const RED = "R"
const BLACK = "B"
const EMPTY = "-"

func initGame(rows int, columns int) game {
	g := game{}
	g.rows = rows
	g.columns = columns

	for i := 0; i < rows*columns; i++ {
		g.board = append(g.board, EMPTY)
	}
	return g
}

func (g game) printGame() {
	fmt.Println()
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.columns; j++ {
			fmt.Print(g.board[g.columns*i+j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g game) playMove(column int, move string) error {

	for i := g.rows - 1; i >= 0; i-- {
		if g.board[g.columns*i+column] == EMPTY {
			g.board[g.columns*i+column] = move
			return nil
		}
	}
	return errors.New("invalid move: column selected is full")
}

func (g game) isWinGame() (bool, string) {

	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.columns; j++ {

			if g.board[g.columns*i+j] != EMPTY {
				// need to check below vertically

				// need to check to the right horizontally

				// need to check to the right below diagonally

				// need to check to the left below diagonally
			}
		}
	}
}
