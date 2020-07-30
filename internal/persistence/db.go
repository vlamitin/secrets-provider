package persistence

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var _db *sql.DB = nil

const DbFileName = "secrets.db"

func Open() *sql.DB {
	db, err := sql.Open("sqlite3", DbFileName)
	if err != nil {
		log.Fatalf("failed to open file %s %v", DbFileName, err)
	}

	_db = db

	return db
}

func PrepareDb(db *sql.DB) {
	prepareSttmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS secrets (id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, key TEXT NOT NULL UNIQUE, value TEXT NOT NULL)")
	if err != nil {
		log.Fatalf("failed to prepare create table statement: %v", err)
	}
	defer prepareSttmt.Close()
	prepareSttmt.Exec()
}

func InsertOne(db *sql.DB, key string, value string) {
	insertSttmt, err := db.Prepare("INSERT INTO secrets (key,value) VALUES (?,?)")
	if err != nil {
		log.Fatalf("failed to prepare insert statement %v", err)
	}
	defer insertSttmt.Close()

	_, insertErr := insertSttmt.Exec(key, value)
	if insertErr != nil {
		fmt.Printf("err when insert key %s: %v\n", key, insertErr)
	}
}

func DeleteOne(db *sql.DB, key string) {
	deleteSttmt, err := db.Prepare("DELETE FROM secrets WHERE key=?")
	if err != nil {
		log.Fatalf("failed to prepare delete statement %v", err)
	}
	defer deleteSttmt.Close()

	_, deleteErr := deleteSttmt.Exec(key)
	if deleteErr != nil {
		fmt.Printf("err when delete key %s: %v\n", key, deleteErr)
	}
}

type SecretRow struct {
	id    int
	key   string
	value string
}

func GetAll(db *sql.DB) []SecretRow {
	rows, err := db.Query("SELECT * FROM secrets")
	if err != nil {
		log.Fatalf("failed to prepare query statement %v", err)
	}
	defer rows.Close()
	res := make([]SecretRow, 20)

	for rows.Next() {
		var id int
		var key string
		var value string

		rows.Scan(&id, &key, &value)
		res = append(res, SecretRow{
			id:    id,
			key:   key,
			value: value,
		})
	}

	return res
}
