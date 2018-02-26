package main

import (
	"connect4/connect4lib"
	"net/http"
)


func main() {

	http.HandleFunc("/ws", connect4lib.ConnectionHandler)
	http.ListenAndServe(":1200", nil)
}
