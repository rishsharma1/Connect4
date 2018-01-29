package main

func main() {
	g := initGame(3, 4)
	g.printGame()
	g.playMove(0, 1, RED)
	g.printGame()
}
