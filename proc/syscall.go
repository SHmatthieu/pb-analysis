package proc

import (
	"database/sql"
	"fmt"
	"time"
)

// the syscall struct
type Syscall struct {
	name         string
	creationDate time.Time
	nTotal       int
	nMax         int
	nMin         int
	mean         float64
	dif          int
	_inter       int
}

func (s *Syscall) String() string {
	return fmt.Sprintf("%s : %d", s.name, s.nTotal)
}

// write the syscall in the db
func (s *Syscall) write(db *sql.DB, uuid string) error {
	statement, err := db.Prepare("INSERT INTO syscall VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(
		uuid,
		s.name,
		s.nTotal,
		s.nMax,
		s.nMin,
		s.mean)
	if err != nil {
		fmt.Println(s.String() + " : " + uuid)
		fmt.Println(err)
		return err
	}

	return nil
}

// CreateSycallTable create syscall table in the db
func CreateSyscallTable(db *sql.DB) error {
	statement, err := db.Prepare(`CREATE TABLE syscall (
		uuid TEXT,
		name TEXT,
		nTotal INT,
		nMax INT,
		nMin INT,
		mean FLOAT,
		PRIMARY KEY (name, uuid));`)
	if err != nil {
		return err
	}
	statement.Exec()
	return nil
}
