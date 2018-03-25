package connect4lib

import (
	"fmt"
)

// LogError outputs the error to stdout but does not exit
// program
func LogError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
