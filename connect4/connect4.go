package connect4

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type game struct {
	board   []string
	columns int
	rows    int
}

// RED player indicator
const RED = "R"

// BLACK player indicator
const BLACK = "B"

// EMPTY board indicator
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
				color := g.board[g.columns*i+j]
				// need to check below vertically
				wonVertically := g.isVertical(i, j, color)
				// need to check to the right horizontally
				wonHorizontally := g.isHorizontal(i, j, color)
				// need to check to the right below diagonally
				wonRightDiagonally := g.isRightDiagonal(i, j, color)
				// need to check to the left below diagonally
				wonLeftDiagonally := g.isLeftDiagonal(i, j, color)

				if wonVertically || wonHorizontally ||
					wonRightDiagonally || wonLeftDiagonally {
					return true, color
				}

			}
		}
	}
	return false, ""
}

func (g game) isVertical(row int, column int, color string) bool {

	if g.rows-row < 4 {
		return false
	}
	for k := 0; k < 4; k++ {
		if g.board[g.columns*(k+row)+column] != color {
			return false
		}
	}
	return true
}

func (g game) isHorizontal(row int, column int, color string) bool {

	if g.columns-column < 4 {
		return false
	}
	for k := 0; k < 4; k++ {
		if g.board[g.columns*row+(column+k)] != color {
			return false
		}
	}
	return true
}

func (g game) isRightDiagonal(row int, column int, color string) bool {

	if g.columns-column < 4 {
		return false
	}
	if g.rows-row < 4 {
		return false
	}
	for k := 0; k < 4; k++ {
		if g.board[g.columns*(k+row)+(column+k)] != color {
			return false
		}
	}
	return true
}

func (g game) isLeftDiagonal(row int, column int, color string) bool {

	if g.rows-row < 4 {
		return false
	}
	if column+1 < 4 {
		return false
	}
	for k := 0; k < 4; k++ {
		if g.board[g.columns*(row+k)+(column-k)] != color {
			return false
		}
	}
	return true
}

func loadGame(filename string) game {
	g := game{}
	b, err := ioutil.ReadFile(filename)
	handleError(err)

	s := strings.Split(string(b), "\n")
	gameInfo := strings.Split(s[0], ",")
	rows, err := strconv.Atoi(gameInfo[0])
	handleError(err)
	columns, err := strconv.Atoi(gameInfo[1])
	handleError(err)
	g.rows = rows
	g.columns = columns

	for i := 1; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			g.board = append(g.board, string(s[i][j]))
		}
	}

	return g

}
