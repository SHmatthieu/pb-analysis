package proc

import (
	"app/tap"
	"os"
	"strconv"
	"sync"
)

// GetAllPids read all folder name in /proc
// return all pid
func GetAllPids() []int {
	data, _ := os.ReadDir("/proc")
	res := make([]int, 1)
	total := 0
	for _, i := range data {
		if isInt([]byte(i.Name())) {
			val, err := strconv.Atoi(i.Name())
			if err == nil {
				res = append(res, val)
				total++
			}
		}
	}
	return res[:total]
}

// init create procs map
func Init() {
	procs = make(map[int]*Proc)
	procsMutex = sync.Mutex{}
}

// Update procs map with new process
func UpdateProcs(pids []int) {
	procsMutex.Lock()
	defer procsMutex.Unlock()
	for _, pid := range pids {
		_, exist := procs[pid]
		if !exist {
			proc, err := LoadProc(pid)
			if err == nil {
				procs[proc.pid] = proc
			} else {

			}
		}
	}

}

// NewSyscallMade add a syscall to a processus
func NewSysCallMade(syscall *tap.Syscall) {
	procsMutex.Lock()
	p, exist := procs[syscall.Pid()]
	procsMutex.Unlock()

	if exist {
		p.NewSyscall(syscall)
	} else if MODE == 0 {
		pids := GetAllPids()
		UpdateProcs(pids)
	}
}
