package proc

import (
	"fmt"
	"strings"
	"time"
)

type procLog struct {
	proc *Proc
	time time.Time
}

var (
	logChan chan string
)

// NeLog create a log struct for a proc
func NewLog(p *Proc) *procLog {
	return &procLog{proc: p, time: time.Now()}
}

// format transfom a log struct into a string to be logged
func (log *procLog) format(syscall string, n int) string {
	var builder strings.Builder
	builder.WriteString(log.time.Format(time.UnixDate))
	builder.WriteString(";")
	builder.WriteString(fmt.Sprintf("%d", log.proc.Pid()))
	builder.WriteString(";")
	builder.WriteString(log.proc.Exe())
	builder.WriteString(";")
	builder.WriteString(syscall)
	builder.WriteString(";")
	builder.WriteString(fmt.Sprintf("%d", n))
	builder.WriteString("\n")
	return builder.String()
}

// write logs data form a log struct in the c channel
func (log *procLog) write(c chan string) {

	//currently only log syscall mades by every process
	for name, syscall := range log.proc.Syscalls() {
		c <- log.format(name, syscall.nTotal)
	}

}

// LogInit init and return the log channel where log are written
func LogInit() chan string {
	logChan = make(chan string)
	log()
	return logChan
}

//log logs data from all proc struct every 2 sec
func log() {
	for _, proc := range procs {
		log := NewLog(proc)
		proc.mutex.Lock()
		log.write(logChan)
		proc.mutex.Unlock()
	}
	time.AfterFunc(time.Second*2, log)
}
