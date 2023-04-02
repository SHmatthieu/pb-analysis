package tap

import (
	"strconv"
	"strings"
)

// syscall struct made to contain data from systemtap
type Syscall struct {
	name string
	pid  int
}

// ParseTap parse byte slice from systemtap
// return a syscall struct made from the data
func ParseTap(data []byte) *Syscall {
	strs := strings.Split(string(data), " ")
	pid, _ := strconv.Atoi(strings.Trim(strs[1], "\n"))
	return &Syscall{name: strs[0], pid: pid}
}

func (s *Syscall) Name() string {
	return s.name
}

func (s *Syscall) Pid() int {
	return s.pid
}
