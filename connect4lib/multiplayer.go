package connect4lib

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// Player represents the connect4 player
// storing their connection and username
type Player struct {
	Conn     net.Conn
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

// Response is the struct that sent by
// the client
type Response struct {
	Action  string
	Content map[string]string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// StartGame starts the connect4 game between playerOne and PlayerTwo with
// rows and columns defined
func StartGame(playerOne Player, playerTwo Player, rows int, columns int) {

	defer playerOne.Conn.Close()
	defer playerTwo.Conn.Close()

	connMap := make(map[string]net.Conn)
	connMap[playerOne.UserName] = playerOne.Conn
	connMap[playerTwo.UserName] = playerTwo.Conn
	g := OnlineGame{}
	g.PlayerColors = make(map[string]string)
	g.OGame = InitGame(rows, columns)

	// Make this allocation random
	g.PlayerColors[playerOne.UserName] = RED
	g.PlayerColors[playerTwo.UserName] = BLACK
	g.CurrentTurn = playerOne.UserName

	// game can now start
	for {

		updateGameMessage(g, playerOne)
		updateGameMessage(g, playerTwo)

		currTurnConn := connMap[g.CurrentTurn]
		request := make([]byte, 512)
		readLen, err := currTurnConn.Read(request)
		HandleError(err)
		response := DecodeResponse(request[:readLen])

		if response.Action == "playMove" {
			err := playMoveHandler(&g, response)
			// need to handle specific wrong move error
			HandleError(err)
			updateGameMessage(g, playerOne)
			updateGameMessage(g, playerTwo)
		}

	}

}

func updateGameMessage(og OnlineGame, player Player) {

	b, _ := json.Marshal(og)
	fmt.Println(string(b))
	sendMessage(string(b), player.Conn)
	//go fmt.Println(string(b))
}

func sendMessage(message string, conn net.Conn) {
	conn.Write([]byte(message))
}

// DecodeResponse decodes the response into a Response struct
func DecodeResponse(data []byte) Response {
	var r Response
	err := json.Unmarshal(data, &r)
	HandleError(err)
	return r
}

// InitPlayerHandler handles requests related to initialising a player
func InitPlayerHandler(conn net.Conn, response Response) Player {

	p := Player{}
	p.Conn = conn
	p.UserName = response.Content["userName"]
	return p

}

// playMoveHandler handles playMove requests
func playMoveHandler(game *OnlineGame, response Response) error {

	err := game.playMove(response)
	HandleError(err)

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
	checkError(err)
	go func(conn *websocket.Conn) {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			fmt.Println(string(msg))
		}
	}(conn)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
