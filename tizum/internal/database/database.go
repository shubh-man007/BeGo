package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DBLite struct {
	DB *sql.DB
}

func InitDB(filepath string) (*DBLite, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		created_at TEXT NOT NULL,
		status INTEGER NOT NULL
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return nil, err
	}

	return &DBLite{DB: db}, nil
}
