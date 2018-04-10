package connect4lib

import (
	"encoding/json"

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

// UpdateGameMessage sends an update game message to the player
func UpdateGameMessage(og OnlineGame, player Player) {

	resp := UpdateResponse{}
	resp.Action = UPDATE
	resp.Og = og
	b, _ := json.Marshal(resp)
	sendMessage(string(b), player.Conn)
}

// InvalidMoveMessage send an invalid move message to the player
// who made an invalid move
func InvalidMoveMessage(player Player) {

	resp := NewResponse()
	resp.Action = INVALIDMOVE
	resp.Content["UserName"] = player.UserName
	b, _ := json.Marshal(resp)
	sendMessage(string(b), player.Conn)

}

// SendMessage sends message to the conn
func sendMessage(message string, conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// SendUpdateGameMessage sends an update game message to the
// waiting player
func SendUpdateGameMessage(player Player, response Response) {

	gameKey := response.Content["GameKey"]
	gameChannel, err := GetGameChannel(gameKey)
	LogError(err)
	og := <-gameChannel
	UpdateGameMessage(*og, player)

}
