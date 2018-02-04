package connect4lib

import (
	"log"
)

// HandleError is used to handle errors
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
