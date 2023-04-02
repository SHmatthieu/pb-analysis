package proc

import (
	"fmt"
	"time"
)

// the prog struct
// represent an application
type Prog struct {
	uuids    []string
	data     map[string]int
	fullPath string
	syscalls map[string]*Syscall
}

var (
	Progs map[string]*Prog
	apps  []string
)

func (p *Prog) String() string {

	return fmt.Sprintf("%s : threads : %d, syscalls : %d, memory : %d", p.fullPath, p.data["nMaxThread"], len(p.syscalls), p.data["memory"])
}

func (p *Prog) Uuid() []string {
	return p.uuids
}

// AddSysCall add a syscall to a prog
// create a new one if it not exist already
func (p *Prog) AddSysCall(name string, nTotal, nMax, nMin int, mean float64) {
	//à ameliorer

	if val, isIn := p.syscalls[name]; isIn {
		if val.nMax < nMax {
			val.nMax = nMax
		}
		if val.nMin > nMin {
			val.nMin = nMin
		}
		newmean := val.mean*float64(val.nTotal) + mean*float64(nTotal)
		val.nTotal += nTotal
		newmean /= float64(val.nTotal)
		val.mean = newmean

	} else {
		p.syscalls[name] = &Syscall{name: name,
			creationDate: time.Now(),
			nTotal:       nTotal,
			nMax:         nMax,
			nMin:         nMin,
			mean:         mean,
			_inter:       0}
	}

}

// CreateProg create a new prog or update data if prog already exist
func CreateProg(uuid, fullPath string, nMaxThread, memory int) error {
	if val, isIn := Progs[fullPath]; isIn {
		val.uuids = append(val.uuids, uuid)
		if val.data["nMaxThread"] < nMaxThread {
			val.data["nMaxThread"] = nMaxThread
		}
		if val.data["memory"] < memory {
			val.data["memory"] = memory
		}
	} else {
		uuids := make([]string, 0)
		uuids = append(uuids, uuid)
		data := make(map[string]int)
		data["nMaxThread"] = nMaxThread
		data["memory"] = memory

		Progs[fullPath] = &Prog{uuids: uuids,
			fullPath: fullPath,
			data:     data,
			syscalls: make(map[string]*Syscall),
		}
	}
	return nil
}

// create progs map with all application in it
func InitProg(a []string) {
	Progs = make(map[string]*Prog)
	apps = a
}

// CompoareProgAndProg it s the core of the monitoring part
func CompareProgAndProc() {
	procsMutex.Lock()
	defer procsMutex.Unlock()
	for _, p := range procs {
		if len(apps) > 0 {
			isIn := false
			for _, app := range apps {
				if app == (p.exe) {

					isIn = true
				}
			}
			if !isIn {
				continue
			}
		}
		if val, isIn := Progs[p.exe]; isIn {

			fmt.Println("monitoring : " + val.fullPath)
			fmt.Println("---------------")
			fmt.Printf("thread : db : %d, live : %d\n", val.data["nMaxThread"], p.data["nMaxThread"])
			fmt.Printf("memory db : %d, live : %d\n", val.data["memory"], p.data["memory"])

			fmt.Printf("n syscalls : %d\n", len(p.syscalls))
			for _, sc := range p.syscalls {
				_, isIn := val.syscalls[sc.name]
				if isIn {
					fmt.Printf("%s , db : %d, live : %d\n", sc.name, val.syscalls[sc.name].nMax, sc.dif)
				} else {
					fmt.Println(sc.name + " NEW SYSCALL")
				}
			}
			if val.data["nMaxThread"]*2 < p.data["nMaxThread"] {
				fmt.Println("suspicious number of thread")
			}
			if val.data["memory"]*2 < p.data["memory"] {
				fmt.Println("suspicious memory usage")

			}
			fmt.Println("---------------")

		} else {
			fmt.Println("not monitoring : " + p.exe + "/" + p.status["Name"])

		}
	}

}
