package tools

import (
	"app/proc"
	"time"
)

// learn is called when the programm is in learn mode
func Learn() {
	proc.UpdateAll()
	time.AfterFunc(time.Second*2, Learn)
}
