package connect4

import (
	"log"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
