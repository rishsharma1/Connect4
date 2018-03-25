package connect4lib

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
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
	GameKey      string
	Players      []Player
}

// Response is the struct that sent
// and received by the client
type Response struct {
	Action  string            `json:"action"`
	Content map[string]string `json:"content"`
}

// UpdateResponse is the struct that is
// sent to update game
type UpdateResponse struct {
	Action string     `json:"action"`
	Og     OnlineGame `json:"og"`
}

// upgrader responsible for upgrading
// http requests to a websockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// UPDATE is the message identifier for
// an update message
const UPDATE = "UPDATE_MESSAGE"

// UPDATEREQUEST is the messaage identifier
// for a request made by the client asking
// for a game update
const UPDATEREQUEST = "UPDATE_REQUEST"

// INIT is the message identifier for
// a player initialisation message
const INIT = "INIT"

// PLAYMOVE is the message identifier for
// a player move message
const PLAYMOVE = "PLAY_MOVE"

// queue of players
var queue = []Player{}
var gameMap = make(map[string]OnlineGame)
var channelMap = make(map[string]chan *OnlineGame)

// InitOnlineGame starts the connect4 game between playerOne and PlayerTwo with
// rows and columns defined
func InitOnlineGame(playerOne Player, playerTwo Player, rows int, columns int) {

	g := OnlineGame{}
	g.PlayerColors = make(map[string]string)
	g.OGame = InitGame(rows, columns)

	// Make this allocation random
	g.PlayerColors[playerOne.UserName] = RED
	g.PlayerColors[playerTwo.UserName] = BLACK
	g.CurrentTurn = playerOne.UserName

	gameKey, _ := exec.Command("uuidgen").Output()
	g.GameKey = string(gameKey)
	gameMap[string(gameKey)] = g
	channelMap[string(gameKey)] = make(chan *OnlineGame)
	g.Players = append(g.Players, playerOne, playerTwo)

	// send inital update game message
	updateGameMessage(g, playerOne)
	updateGameMessage(g, playerTwo)

}

// updateGameMessage sends an update game message to the player
func updateGameMessage(og OnlineGame, player Player) {

	resp := UpdateResponse{}
	resp.Action = UPDATE
	resp.Og = og
	b, _ := json.Marshal(resp)
	sendMessage(string(b), player.Conn)
}

// sendMessage sends message to the conn
func sendMessage(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// sendUpdateGameMessage sends an update game message to the
// waiting player
func sendUpdateGameMessage(player Player, response Response) {

	og := <-channelMap[response.Content["GameKey"]]
	updateGameMessage(*og, player)

}

// DecodeResponse decodes the response into a Response struct
func DecodeResponse(data []byte) Response {
	var r Response
	err := json.Unmarshal(data, &r)
	LogError(err)
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
	LogError(err)
	userName := response.Content["UserName"]
	for player := range game.PlayerColors {
		if userName != player {
			game.CurrentTurn = player
			break
		}
	}

	wonGame, winner := game.OGame.IsWinGame()
	if wonGame {
		game.Winner = winner
		game.GameState = "won"

	}
	log.Println(game)
	gameChannel := channelMap[game.GameKey]
	go func() { gameChannel <- game }()
	return err
}

// playMove plays the move specified by the username
func (og OnlineGame) playMove(response Response) error {

	column, err := strconv.Atoi(response.Content["Column"])
	userName := response.Content["UserName"]
	color := og.PlayerColors[userName]
	err = og.OGame.PlayMove(column, color)
	return err
}

// ConnectionHandler handles new connections
func ConnectionHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	LogError(err)

	go func(conn *websocket.Conn) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			r := DecodeResponse(msg)
			if r.Action == INIT {
				player := InitPlayerHandler(conn, r)
				log.Println("Player Connected:", player.UserName)
				queue = append(queue, player)
			} else if r.Action == PLAYMOVE {
				g := gameMap[r.Content["GameKey"]]
				err := playMoveHandler(&g, r)
				LogError(err)
				player := InitPlayerHandler(conn, r)
				updateGameMessage(g, player)
			} else if r.Action == UPDATEREQUEST {
				player := InitPlayerHandler(conn, r)
				sendUpdateGameMessage(player, r)
			}
			if len(queue) > 1 {
				p1 := queue[0]
				p2 := queue[1]
				queue = queue[2:]
				log.Println("Starting Game:", p1.UserName, "vs", p2.UserName)
				InitOnlineGame(p1, p2, 6, 7)
			}
		}
	}(conn)
}
