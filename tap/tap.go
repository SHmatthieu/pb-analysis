package tap

/*
 this package contain
 functions and structures
 to interact with data from systemtap

 P3 he-arc
 matthieu barbot 2021/2022
*/

import (
	"bufio"
	"fmt"
	"os"
)

// ReadStdin read standard input
func ReadStdin(c chan []byte) {
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			fmt.Print(err.Error())
			close(c)
			break
		}
		c <- data
	}
}
