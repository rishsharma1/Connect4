package connect4lib

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// InitPlayerHandler handles requests related to initialising a player
func InitPlayerHandler(conn *websocket.Conn, response Response) Player {

	p := Player{}
	p.Conn = conn
	p.UserName = response.Content["UserName"]
	return p

}

// PlayMoveHandler handles playMove requests
func PlayMoveHandler(conn *websocket.Conn, response Response) error {

	gameKey := response.Content["GameKey"]
	game := GetOnlineGame(gameKey)
	player := InitPlayerHandler(conn, response)
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
		gameChannel, err := GetGameChannel(gameKey)
		LogError(err)
		log.Println(gameChannel)

		go func(gameChannel chan *OnlineGame, game *OnlineGame) {
			gameChannel <- game
		}(gameChannel, game)
		UpdateGameMessage(*game, player)

	} else {

		LogError(err)
		InvalidMoveMessage(player)

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
			action := r.Action
			if action == INIT {
				p := InitPlayerHandler(conn, r)
				InsertPlayerQueue(p)
				log.Println("Player Connected:", p.UserName)
			} else if action == PLAYMOVE {
				PlayMoveHandler(conn, r)
			} else if action == UPDATEREQUEST {
				player := InitPlayerHandler(conn, r)
				SendUpdateGameMessage(player, r)
			}
			if GetPlayerQueueLen() > 1 {
				p1, err := PopPlayerQueue()
				LogError(err)
				p2, err := PopPlayerQueue()
				LogError(err)
				log.Println("Starting Game:", p1.UserName, "vs", p2.UserName)
				InitOnlineGame(p1, p2, 6, 7)
			}
		}
	}(conn)
}
