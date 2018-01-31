package connect4

import (
	"testing"
)

func TestVerticalWinState(t *testing.T) {

	testName := "tests/vertical_test"
	g := loadGame(testName)
	winGame, _ := g.isWinGame()
	if !winGame {
		t.Errorf("Expected a true vertical win state, got false")
	}
}

func TestHorizontalWinState(t *testing.T) {

	testName := "tests/horizontal_test"
	g := loadGame(testName)
	winGame, _ := g.isWinGame()
	if !winGame {
		t.Errorf("Expected a true horizontal win state, got false")
	}
}

func TestRightDigonalWinState(t *testing.T) {

	testName := "tests/right_diagonal_test"
	g := loadGame(testName)
	winGame, _ := g.isWinGame()
	if !winGame {
		t.Errorf("Expected a true right diagonal win state, got false")
	}
}

func TestLeftDiagonalWinState(t *testing.T) {

	testName := "tests/left_diagonal_test"
	g := loadGame(testName)
	winGame, _ := g.isWinGame()
	if !winGame {
		t.Errorf("Expected a true left diagonal win state, got false")
	}
}
