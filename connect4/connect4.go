package main

import (
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

func (g game) playMove(row int, column int, move string) {
	g.board[g.columns*row+column] = move
}
