package connect4lib

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Player represents the connect4 player
// storing their connection and username
type Player struct {
	Conn     *websocket.Conn
	UserName string
}

// OnlineGame is a struct that represents the
// the state of the game
type OnlineGame struct {
	OGame        Game
	Moves        int
	PlayerColors map[string]string
	CurrentTurn  string
	GameState    string
	Winner       string
}

// Response is the struct that sent
// and received by the client
type Response struct {
	Action  string            `json:"action"`
	Content map[string]string `json:"content"`
}

// upgrader responsible for upgrading
// http requests to a websockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// queue of players
var queue = []Player{}

// StartGame starts the connect4 game between playerOne and PlayerTwo with
// rows and columns defined
func StartGame(playerOne Player, playerTwo Player, rows int, columns int) {

	defer playerOne.Conn.Close()
	defer playerTwo.Conn.Close()

	connMap := make(map[string]*websocket.Conn)
	connMap[playerOne.UserName] = playerOne.Conn
	connMap[playerTwo.UserName] = playerTwo.Conn
	g := OnlineGame{}
	g.PlayerColors = make(map[string]string)
	g.OGame = InitGame(rows, columns)

	// Make this allocation random
	g.PlayerColors[playerOne.UserName] = RED
	g.PlayerColors[playerTwo.UserName] = BLACK
	g.CurrentTurn = playerOne.UserName

	// send inital update game message
	updateGameMessage(g, playerOne)
	updateGameMessage(g, playerTwo)

	// game can now start
	for {

		currTurnConn := connMap[g.CurrentTurn]
		_, msg, err := currTurnConn.ReadMessage()
		HandleError(err)
		response := DecodeResponse(msg)

		if response.Action == "playMove" {
			err := playMoveHandler(&g, response)
			// need to handle specific wrong move error
			HandleError(err)
			updateGameMessage(g, playerOne)
			updateGameMessage(g, playerTwo)
		}

	}

}

// updateGameMessage sends an update game message to the player
func updateGameMessage(og OnlineGame, player Player) {

	b, _ := json.Marshal(og)
	sendMessage(string(b), player.Conn)
}

// sendMessage sends message to the conn
func sendMessage(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// DecodeResponse decodes the response into a Response struct
func DecodeResponse(data []byte) Response {
	var r Response
	err := json.Unmarshal(data, &r)
	HandleError(err)
	return r
}

// InitPlayerHandler handles requests related to initialising a player
func InitPlayerHandler(conn *websocket.Conn, response Response) Player {

	p := Player{}
	p.Conn = conn
	p.UserName = response.Content["UserName"]
	return p

}

// playMoveHandler handles playMove requests
func playMoveHandler(game *OnlineGame, response Response) error {

	err := game.playMove(response)
	CheckError(err)

	wonGame, winner := game.OGame.IsWinGame()
	if wonGame {
		game.Winner = winner
		game.GameState = "won"

	}
	return err
}

// playMove plays the move specified by the username
func (og OnlineGame) playMove(response Response) error {

	column, err := strconv.Atoi(response.Content["column"])
	HandleError(err)
	color := og.PlayerColors[og.CurrentTurn]
	err = og.OGame.PlayMove(column, color)
	return err
}

// ConnectionHandler handles new connections
func ConnectionHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	CheckError(err)

	go func(conn *websocket.Conn) {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			r := DecodeResponse(msg)
			if r.Action == "init" {
				player := InitPlayerHandler(conn, r)
				log.Println("Player Connected:", player.UserName)
				queue = append(queue, player)
			}
			if len(queue) > 1 {
				p1 := queue[0]
				p2 := queue[1]
				log.Println("Starting Game:", p1.UserName, "vs", p2.UserName)
				go StartGame(p1, p2, 6, 7)
			}
		}
	}(conn)
}
