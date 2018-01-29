package main

import (
	"log"
)

func main() {
	g := initGame(3, 4)
	g.printGame()
	g.playMove(0, RED)
	g.playMove(0, RED)
	g.playMove(0, RED)
	err := g.playMove(0, RED)
	if err != nil {
		log.Fatal(err)
	}
	g.printGame()
}
