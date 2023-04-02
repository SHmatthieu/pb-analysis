package proc

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// remove return if a char is space, a tab, an LF or not
func remove(char int32) bool {
	if char == 9 || char == 32 || char == 10 {
		return true
	}
	return false
}

// isInit return if every byte in a slice of byte is an Integer
func isInt(data []byte) bool {
	for _, d := range data {
		if d < 48 || d > 57 {
			return false
		}
	}
	return true
}

// parseStatus parse the Status file in /proc/{pid}/
// return a map with {"name":value,...}
func parseStatus(f *os.File) map[string]string {
	scanner := bufio.NewScanner(f)
	result := make(map[string]string)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ":")
		if len(data) == 2 {
			result[data[0]] = strings.TrimFunc(data[1], remove)
		}
	}
	return result
}

//parseMem return an memory value in int
func parseMem(data string) int {
	value, err := strconv.ParseInt(data[:len(data)-3], 10, 64)
	if err != nil {
		return 0
	}
	return int(value)
}
