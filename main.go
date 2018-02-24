package main

import (
	"connect4/connect4lib"
	"net/http"
)

/*
func main() {

	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	connect4lib.HandleError(err)
	listner, err := net.ListenTCP("tcp", tcpAddr)
	connect4lib.HandleError(err)
	//queue := []connect4lib.Player{}

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}
		data := make([]byte, 512)
		readLen, err := conn.Read(data)
		fmt.Println(string(data[:readLen]))
*/
/*
	connect4lib.HandleError(err)
	response := connect4lib.DecodeResponse(data[:readLen])
	player := connect4lib.InitPlayerHandler(conn, response)
	fmt.Println(response)

	if len(queue) > 0 {
		// start game
		playerOne := queue[0]
		playerTwo := player
		queue = queue[1:]
		go connect4lib.StartGame(playerOne, playerTwo, 6, 7)
	} else {
		queue = append(queue, player)
	}
*/
/*
	}

}

*/

/*
func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		so.Join("chat")
		log.Println("on connection ", so.Id())

		so.On("subscribeToTimer", func(msg string) {
			log.Println("subscribe to timer with time interval", msg)
			so.Emit("timer", time.Now().String())
		})
		so.On("disconnection", func() {
			log.Println("error:", err)
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:1200...")
	log.Fatal(http.ListenAndServe("localhost:1200", nil))

}
*/

func main() {

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.)
			http.ServeFile(w, r, "index.html")
		})
	*/

	http.HandleFunc("/ws", connect4lib.ConnectionHandler)
	http.ListenAndServe(":1200", nil)
}
