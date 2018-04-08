package connect4lib

import (
	"encoding/json"
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

// QUEUE is queue to manage players waiting to play
var QUEUE = []Player{}

// GAMEMAP is a maps unique game keys to game states
var GAMEMAP = make(map[string]OnlineGame)

// CHANNELMAP maps unique game keys to game channels
var CHANNELMAP = make(map[string]chan *OnlineGame)

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
	GAMEMAP[string(gameKey)] = g
	CHANNELMAP[string(gameKey)] = make(chan *OnlineGame)
	g.Players = append(g.Players, playerOne, playerTwo)

	// send inital update game message
	updateGameMessage(g, playerOne)
	updateGameMessage(g, playerTwo)

}

// GetGameChannel returns a game channel for the given
// game key
func GetGameChannel(gameKey string) chan *OnlineGame {
	return CHANNELMAP[gameKey]
}

// GetOnlineGame returns the state of the online game
// for the given game key
func GetOnlineGame(gameKey string) OnlineGame {
	return GAMEMAP[gameKey]
}

// PopPlayerQueue returns the next Player in the queue
// waiting for a game
func PopPlayerQueue() Player {
	p := QUEUE[0]
	QUEUE = QUEUE[1:]
	return p
}

// InsertPlayerQueue appends the player to the queue of
// players
func InsertPlayerQueue(player Player) {
	QUEUE = append(QUEUE, player)
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

// SendUpdateGameMessage sends an update game message to the
// waiting player
func SendUpdateGameMessage(player Player, response Response) {

	og := <-CHANNELMAP[response.Content["GameKey"]]
	updateGameMessage(*og, player)

}

// DecodeResponse decodes the response into a Response struct
func DecodeResponse(data []byte) Response {
	var r Response
	err := json.Unmarshal(data, &r)
	LogError(err)
	return r
}

// PlayMove plays the move specified by the username
func (og OnlineGame) PlayMove(response Response) error {

	column, err := strconv.Atoi(response.Content["Column"])
	userName := response.Content["UserName"]
	color := og.PlayerColors[userName]
	err = og.OGame.PlayMove(column, color)
	return err
}
