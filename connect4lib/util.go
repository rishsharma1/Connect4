package connect4lib

import (
	"fmt"
	"log"
)

// HandleError is used to handle errors
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Outputs the error to stdout but does not exit
// program
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
