package models

import "time"

type Task struct {
	ID        int       `json:"id"`
	TaskName  string    `json:"task"`
	CreatedAt time.Time `json:"created_at"`
	Status    bool      `json:"status"`
}
