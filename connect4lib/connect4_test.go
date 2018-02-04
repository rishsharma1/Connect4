package connect4lib

import (
	"testing"
)

func TestInitGame(t *testing.T) {
	g := InitGame(5, 4)
	if len(g.board) != 5*4 {
		t.Errorf("expected the length of the game board to be 20, got %d", len(g.board))
	}
}

func TestLoadGame(t *testing.T) {

	testName := "tests/load_game_test"
	g := LoadGame(testName)
	if len(g.board) != g.rows*g.columns {
		t.Errorf("expected the length of the game board to be %d, got %d", g.rows*g.columns, len(g.board))
	}
}

func TestPlayMove(t *testing.T) {

	g := InitGame(4, 4)
	g.PlayMove(0, RED)
	if g.board[g.columns*3] != RED {
		t.Errorf("expected played move to be 'R' at (0,3), got %v", g.board[g.columns*3])
	}
}

func TestVerticalWinState(t *testing.T) {

	testName := "tests/vertical_test"
	g := LoadGame(testName)
	winGame, _ := g.IsWinGame()
	if !winGame {
		t.Errorf("Expected a true vertical win state, got false")
	}
}

func TestHorizontalWinState(t *testing.T) {

	testName := "tests/horizontal_test"
	g := LoadGame(testName)
	winGame, _ := g.IsWinGame()
	if !winGame {
		t.Errorf("Expected a true horizontal win state, got false")
	}
}

func TestRightDigonalWinState(t *testing.T) {

	testName := "tests/right_diagonal_test"
	g := LoadGame(testName)
	winGame, _ := g.IsWinGame()
	if !winGame {
		t.Errorf("Expected a true right diagonal win state, got false")
	}
}

func TestLeftDiagonalWinState(t *testing.T) {

	testName := "tests/left_diagonal_test"
	g := LoadGame(testName)
	winGame, _ := g.IsWinGame()
	if !winGame {
		t.Errorf("Expected a true left diagonal win state, got false")
	}
}
