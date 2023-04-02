package main

import (
	"app/data"
	"app/log"
	"app/proc"
	"app/tap"
	"app/tools"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// entry point of the programm
func main() {
	fmt.Println("p3 !")

	proc.Init()

	pids := loadConfigPid("./config/PID.CONFIG")
	apps := loadConfig("./config/PROG.CONFIG")
	proc.SetModeConfig()

	if len(pids) == 0 {
		proc.SetModeAllPid()
		pids = proc.GetAllPids()

	}
	proc.UpdateProcs(pids)

	go log.Init("./logs.txt", proc.LogInit)

	// createe a chanel to communicate betwin the function that
	// will parse input from Systemtap script and the function that will
	// process it
	c := make(chan []byte, 1)
	go tap.ReadStdin(c)

	mode := getMode()

	db, err := data.InitDB("./test.db", mode)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(-1)
	}

	if mode == 1 {
		SetupCloseHandler(db)
		data.CreateTables(db, proc.CreateProcTable, proc.CreateSyscallTable)
		go tools.Learn()
	} else {
		go tools.Monitor(db, apps)
	}

	// process data from systemTap
	for v := range c {
		syscall := tap.ParseTap(v)
		proc.NewSysCallMade(syscall)
	}

	db.Close()
}

// getMode read the command line argument and
// return the current mode, learning or monitor by default learning
func getMode() int {
	if len(os.Args) > 1 {
		if strings.Compare(os.Args[1], "learn") == 0 || strings.Compare(os.Args[1], "1") == 0 {
			return 1
		}
	}
	return 0
}

// loadConfigPid pase the config file for pids
// return a slice with every pid in the config file
func loadConfigPid(path string) []int {
	pids := loadConfig(path)
	res := make([]int, len(pids))
	for _, p := range pids {
		val, err := strconv.Atoi(p)
		if err == nil {
			res = append(res, val)
		}
	}
	return res
}

// loadConfig open a config file and parse the content
// return a slice with every line in the file
func loadConfig(path string) []string {
	readFile, err := os.Open(path)
	res := make([]string, 0)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		res = append(res, fileScanner.Text())
	}

	readFile.Close()
	return res
}

// SetupCloseHandler catch the interupt signal,
// write data in the DB and close the programm
func SetupCloseHandler(db *sql.DB) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		proc.WriteAll(db)
		os.Exit(0)
	}()
}
