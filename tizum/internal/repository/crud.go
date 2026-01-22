package repository

import (
	"log"
	"time"

	"github.com/shubh-man007/Tizu/internal/database"
	"github.com/shubh-man007/Tizu/internal/models"
)

type TizOrch struct {
	Orch *database.DBLite
}

func NewTiz(db *database.DBLite) *TizOrch {
	return &TizOrch{
		Orch: db,
	}
}

func (t *TizOrch) CreateTask(taskName string) error {
	query := "INSERT INTO tasks (task, created_at, status) VALUES (?, ?, ?)"
	_, err := t.Orch.DB.Exec(query, taskName, time.Now().Format(time.RFC3339), 0)
	return err
}

func (t *TizOrch) ReadTasks() ([]models.Task, error) {
	query := "SELECT id, task, created_at, status FROM tasks"
	rowTasks, err := t.Orch.DB.Query(query)
	if err != nil {
		log.Printf("Could not fetch tasks: %v", err)
		return nil, err
	}
	defer rowTasks.Close()

	var tasks []models.Task

	for rowTasks.Next() {
		var i models.Task
		var createdAtStr string
		var statusInt int

		if err := rowTasks.Scan(
			&i.ID,
			&i.TaskName,
			&createdAtStr,
			&statusInt,
		); err != nil {
			return nil, err
		}

		parsedTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, err
		}
		i.CreatedAt = parsedTime

		i.Status = statusInt == 1

		tasks = append(tasks, i)
	}

	if err := rowTasks.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TizOrch) UpdateTask(id string, done bool) error {
	query := "UPDATE tasks SET status = ? WHERE id = ?"
	status := 0
	if done {
		status = 1
	}

	_, err := t.Orch.DB.Exec(query, status, id)
	return err
}

func (t *TizOrch) DeleteTask(id string) error {
	query := "DELETE FROM tasks WHERE id = ?"

	_, err := t.Orch.DB.Exec(query, id)
	return err
}
