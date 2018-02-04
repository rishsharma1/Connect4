package connect4lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Game stores the structure of the game
type Game struct {
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

// InitGame will initialise a connect4 game
func InitGame(rows int, columns int) Game {
	g := Game{}
	g.rows = rows
	g.columns = columns

	for i := 0; i < rows*columns; i++ {
		g.board = append(g.board, EMPTY)
	}
	return g
}

// PrintGame prints the game to stdout
func (g Game) PrintGame() {
	fmt.Println()
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.columns; j++ {
			fmt.Print(g.board[g.columns*i+j])
		}
		fmt.Println()
	}
	fmt.Println()
}

// PlayMove plays the move at the specified column
func (g Game) PlayMove(column int, move string) error {

	for i := g.rows - 1; i >= 0; i-- {
		if g.board[g.columns*i+column] == EMPTY {
			g.board[g.columns*i+column] = move
			return nil
		}
	}
	return errors.New("invalid move: column selected is full")
}

// IsWinGame checks to see if the game is in a win state
func (g Game) IsWinGame() (bool, string) {

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

// isVertical checks to see if the game contains a vertiacal win
func (g Game) isVertical(row int, column int, color string) bool {

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

// isHorizontal checks to see if the game contains a horizontal win
func (g Game) isHorizontal(row int, column int, color string) bool {

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

// isRightDiagonal checks to see if the game contains a right diagonal win
func (g Game) isRightDiagonal(row int, column int, color string) bool {

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

// isLeftDiagonal checks to see if the game contains a left diagonal win
func (g Game) isLeftDiagonal(row int, column int, color string) bool {

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

// LoadGame loads a game from a file
func LoadGame(filename string) Game {
	g := Game{}
	b, err := ioutil.ReadFile(filename)
	HandleError(err)

	s := strings.Split(string(b), "\n")
	gameInfo := strings.Split(s[0], ",")
	rows, err := strconv.Atoi(gameInfo[0])
	HandleError(err)
	columns, err := strconv.Atoi(gameInfo[1])
	HandleError(err)
	g.rows = rows
	g.columns = columns

	for i := 1; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			g.board = append(g.board, string(s[i][j]))
		}
	}

	return g

}
