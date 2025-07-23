package main

import (
	"log"
)

func main() {
	connStr := "postgres://gobank:gobankpass@localhost:5432/gobank?sslmode=disable"
	pgxDB, err := NewPostgresStorage(connStr)
	if err != nil {
		log.Fatalf("Error in creating DB %s", err.Error())
	}

	if err := pgxDB.CreateAccountsTable(); err != nil {
		log.Fatalf("Error creating accounts table: %s", err.Error())
	}

	srv := NewAPIServer(":8080", pgxDB)
	srv.Run()
}
