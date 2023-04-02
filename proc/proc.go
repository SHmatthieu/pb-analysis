package proc

/*
 this package contain
 functions and structures
 about processus and syscall

 P3 he-arc
 matthieu barbot 2021/2022
*/

import (
	"app/tap"
	"container/list"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// map with all pid and processus corresponding
	procs      map[int]*Proc
	procsMutex sync.Mutex

	// MODE=0 -> all pids
	MODE int = 0
)

// the processus struct, a proc is
// a representation in the code of a process in the system
type Proc struct {
	uuid     string
	pid      int
	status   map[string]string
	exe      string
	syscalls map[string]*Syscall
	data     map[string]int
	mutex    sync.Mutex
	files    *list.List
}

// Test is test function for devloppement purpose
func Test() {
	Init()
	pids := GetAllPids()
	UpdateProcs(pids)
	for _, proc := range procs {
		fmt.Printf("%d %s:  max threads: %d, mem: %s\n", proc.pid, proc.exe, proc.data["nMaxThread"], proc.status["VmRSS"])
		proc.PrintFiles()

	}
}

// LoadPorc take a pid as argument and retun the corresponding
// process struct with data in it from /proc/{pid}/..
func LoadProc(pid int) (*Proc, error) {
	exe, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return nil, err
	}

	filesDesc, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/fd/", pid))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fileList := list.New()
	for _, file := range filesDesc {
		path, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", pid, file.Name()))
		if err == nil {
			fileList.PushBack(path)
		}
	}

	newProc := &Proc{uuid: uuid.New().String(),
		pid:      pid,
		status:   make(map[string]string),
		exe:      exe,
		data:     make(map[string]int),
		syscalls: make(map[string]*Syscall),
		files:    fileList,
		mutex:    sync.Mutex{}}

	newProc.UpdateStatus()
	newProc.data["nMaxThread"] = newProc.Thread()
	newProc.data["memory"] = newProc.Memory()

	return newProc, nil
}

func (p *Proc) String() string {
	data := ""
	for _, s := range p.syscalls {
		data += s.String() + "\n"
	}
	return fmt.Sprintf("pid : %d, uuid : %s, name : %s, exe : %s, max threads: %d, mem: %s, syscalls : %s",
		p.pid,
		p.uuid,
		p.status["Name"],
		p.exe,
		p.data["nMaxThread"],
		p.status["VmRSS"],
		data)
}

// IsRunning return is a process is still runing or not
func (p *Proc) IsRunning() bool {
	f, err := os.Open(fmt.Sprintf("/proc/%d/status", p.pid))
	if err != nil {
		return false
	}
	f.Close()
	return true
}

// CreateProcTable create the proc table in the data base
func CreateProcTable(db *sql.DB) error {
	statement, err := db.Prepare(`CREATE TABLE proc (
		uuid TEXT,
		name TEXT,
		pathToExe TEXT, 
		nMaxThread INT,
		memory INT,
		PRIMARY KEY (name, pathToExe, uuid));`)
	if err != nil {
		return err
	}
	statement.Exec()
	return nil
}

// write data from the proc struct in the database
func (p *Proc) write(db *sql.DB) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	statement, err := db.Prepare("INSERT INTO proc VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(p.uuid,
		p.status["Name"],
		p.exe,
		p.data["nMaxThread"],
		p.data["memory"])
	if err != nil {
		return err
	}
	return nil
}

// write all procs the procs map in the database
func WriteAll(db *sql.DB) {
	for _, p := range procs {
		err := p.write(db)
		for _, syscall := range p.syscalls {
			syscall.write(db, p.uuid)
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// NewSycall create and add a new syscall
// to the process struct.
func (p *Proc) NewSyscall(syscall *tap.Syscall) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_, exist := p.syscalls[syscall.Name()]
	if !exist {
		p.syscalls[syscall.Name()] = &Syscall{
			name:         syscall.Name(),
			creationDate: time.Now(),
			nTotal:       1,
			nMax:         1,
			nMin:         999999999,
			mean:         0,
			dif:          0,
			_inter:       0}
		return
	}
	p.syscalls[syscall.Name()].nTotal++
}

func (p *Proc) Pid() int {
	return p.pid
}

func (p *Proc) Exe() string {
	return p.exe
}

func (p *Proc) Syscalls() map[string]*Syscall {
	return p.syscalls
}

// UpdateSyscalls update all syscalls struct
// of this call
func (p *Proc) UpdateSyscalls() {
	p.mutex.Lock()
	for _, syscall := range p.syscalls {
		syscall.dif = syscall.nTotal - syscall._inter
		syscall._inter = syscall.nTotal
		syscall.mean = float64(syscall.nTotal) / float64(time.Since(syscall.creationDate).Seconds())
		if syscall.dif > syscall.nMax {
			syscall.nMax = syscall.dif
		}
		if syscall.dif < syscall.nMin {
			syscall.nMin = syscall.dif
		}
	}
	p.mutex.Unlock()
}

// update all proc struct in procs with fresh data
func UpdateAll() {
	procsMutex.Lock()
	defer procsMutex.Unlock()

	for _, proc := range procs {
		//fmt.Println(proc)
		proc.UpdateSyscalls()
		proc.UpdateStatus()
		proc.UpdateData()
	}
}

// Thread return the number of thread of this proc
func (p *Proc) Thread() int {
	if thread, exist := p.status["Threads"]; exist {
		value, err := strconv.ParseInt(thread, 10, 64)
		if err != nil {
			return 0
		}
		return int(value)
	}
	return 0
}

func (p *Proc) Memory() int {
	if memory, exist := p.status["VmRSS"]; exist {
		return parseMem(memory)
	}
	return 0
}

// UpdataStatus Update data from the status file
// of this proc
func (p *Proc) UpdateStatus() error {
	f, err := os.Open(fmt.Sprintf("/proc/%d/status", p.pid))
	if err != nil {
		return err
	}

	p.status = parseStatus(f)
	f.Close()
	return nil

}

// UpdateData update the data map of proc
// with data parsed from the data file
func (p *Proc) UpdateData() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Thread() > p.data["nMaxThread"] {
		p.data["nMaxThread"] = p.Thread()
	}
	if p.Memory() > p.data["memory"] {
		p.data["memory"] = p.Memory()
	}
}

// PrintFiles print all files open by a proc
func (p *Proc) PrintFiles() {
	for f := p.files.Front(); f != nil; f = f.Next() {
		fmt.Println(f)
	}
}

func SetModeAllPid() {
	MODE = 0
}
func SetModeConfig() {
	MODE = 1
}
