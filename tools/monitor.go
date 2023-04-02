package tools

import (
	"app/proc"
	"database/sql"
	"fmt"
	"time"
)

// Monior function is called when the programm is in monitor mod
func Monitor(db *sql.DB, apps []string) {

	// a séparer en plusieurs fonction

	// read data from the DB
	stmt, err := db.Prepare("select name, pathToExe, uuid, nMaxThread, memory from proc")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	proc.InitProg(apps)

	for rows.Next() {
		name := ""
		pathToExe := ""
		uuid := ""
		nMaxThread := 0
		memory := 0
		rows.Scan(&name, &pathToExe, &uuid, &nMaxThread, &memory)
		//fmt.Println(name + " " + pathToExe + " " + uuid + " " + string(nMaxThread) + " " + string(memory))
		proc.CreateProg(uuid, pathToExe, nMaxThread, memory)
		if err = rows.Err(); err != nil {
			fmt.Println(err)
		}
	}

	stmt2, err := db.Prepare("select name, nTotal, nMax, nMin, mean from syscall where uuid = ?")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt2.Close()

	for _, prog := range proc.Progs {
		for _, uuid := range prog.Uuid() {

			rows2, err := stmt2.Query(uuid)
			if err != nil {
				fmt.Println(err)
			}
			for rows2.Next() {

				name := ""
				nTotal := 0
				nMax := 0
				nMin := 0
				mean := 0.0
				rows2.Scan(&name, &nTotal, &nMax, &nMin, &mean)
				prog.AddSysCall(name, nTotal, nMax, nMin, mean)

			}

			rows2.Close()
		}

	}

	for _, prog := range proc.Progs {
		fmt.Println(prog)
	}

	fmt.Println("START MONITORING...")
	RecMonitor()

}

func RecMonitor() {
	proc.UpdateAll()
	proc.CompareProgAndProc()
	time.AfterFunc(time.Second*2, RecMonitor)

}
