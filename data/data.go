package data

/*
 this package contain
 functions and structures
 to interact with the database

 P3 he-arc
 matthieu barbot 2021/2022
*/

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB open or create the database
// return the sql.DB instance
func InitDB(path string, mode int) (*sql.DB, error) {
	//a corriger (séparer en 2 fonction)

	_, err := os.Open(path)
	var new bool = false
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		new = true
	}

	if mode == 0 && new {
		return nil, errors.New("no DB")
	}
	if mode == 1 && !new {
		return nil, errors.New("DB already created")
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if mode == 0 && !validDB(db) {
		return nil, errors.New("error in DB")
	}
	return db, nil
}

func validDB(db *sql.DB) bool {
	return true
}

// CreateTables create the tables in the db
// take as arugments the db instance and functions that will create a table
func CreateTables(db *sql.DB, creators ...func(*sql.DB) error) error {
	for _, creator := range creators {
		err := creator(db)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}
