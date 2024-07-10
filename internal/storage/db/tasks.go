package db

import (
	_ "embed"
	"time"

	"github.com/google/uuid"

	"go-effective-mobile/internal/models"
)

//go:embed queries/insert_task.sql
var insertTask string

func NewTask(t *models.Task) (int, error) {
	var taskID int

	if t.Title == "" {
		t.Title = uuid.New().String()
	}

	err := client.Pool.QueryRow(client.Ctx, insertTask, t.UserID, t.Title).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

//go:embed queries/get_task.sql
var getTask string

func GetTask(id int, at *time.Time) (*models.Task, error) {
	var task *models.Task

	err := client.Pool.QueryRow(client.Ctx, getTask, id, &at).Scan(&task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

//go:embed queries/start_task.sql
var startTask string

func StartTask(id int) error {
	_, err := client.Pool.Exec(client.Ctx, startTask, id)
	return err
}

//go:embed queries/stop_task.sql
var stopTask string

func StopTask(id int) (string, error) {
	var duration string

	err := client.Pool.QueryRow(client.Ctx, stopTask, id).Scan(&duration)
	if err != nil {
		return "", err
	}

	return duration, nil
}
