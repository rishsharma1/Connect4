package connect4lib

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

// INVALIDMOVE is the message identifier for
// a player move invalid message
const INVALIDMOVE = "INVALID_MOVE"

// InitPlayerHandler handles requests related to initialising a player
func InitPlayerHandler(conn *websocket.Conn, response Response) Player {

	p := Player{}
	p.Conn = conn
	p.UserName = response.Content["UserName"]
	return p

}

// PlayMoveHandler handles playMove requests
func PlayMoveHandler(game *OnlineGame, response Response) error {

	err := game.PlayMove(response)
	if err == nil {

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
		gameChannel := GetGameChannel(game.GameKey)
		log.Println(gameChannel)
		go func() { gameChannel <- game }()
		return err
	}

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
				p := InitPlayerHandler(conn, r)
				InsertPlayerQueue(p)
				log.Println("Player Connected:", p.UserName)
			} else if r.Action == PLAYMOVE {
				gameKey := r.Content["GameKey"]
				g := GetOnlineGame(gameKey)
				err := PlayMoveHandler(&g, r)
				LogError(err)
				player := InitPlayerHandler(conn, r)
				updateGameMessage(g, player)
			} else if r.Action == UPDATEREQUEST {
				player := InitPlayerHandler(conn, r)
				SendUpdateGameMessage(player, r)
			}
			if len(queue) > 1 {
				p1 := PopPlayerQueue()
				p2 := PopPlayerQueue()
				log.Println("Starting Game:", p1.UserName, "vs", p2.UserName)
				InitOnlineGame(p1, p2, 6, 7)
			}
		}
	}(conn)
}
