package main

import (
	"fmt"
	"log"

	"github.com/shubh-man007/Tizu/internal/database"
	"github.com/shubh-man007/Tizu/internal/repository"
)

func main() {
	db, err := database.InitDB("tizu.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	tizOrch := repository.NewTiz(db)

	tasks, err := tizOrch.ReadTasks()
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
	}

	fmt.Println(tasks)
}
