package log

/*
 this package contain
 functions and structures
 to logs data in a text file

 P3 he-arc
 matthieu barbot 2021/2022
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Init create a log file(path) and log data
// from the c channel create by the initFunc
func Init(path string, initFunc func() chan string) {
	fmt.Println("start  logging...")
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	logWriter := bufio.NewWriter(file)

	c := initFunc()
	for v := range c {
		logWriter.WriteString(v)
	}
}
