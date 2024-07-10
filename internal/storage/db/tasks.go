package db

import (
	_ "embed"
	"fmt"
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
	task := &models.Task{}

	err := client.Pool.QueryRow(client.Ctx, getTask, id, at).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.StartedAt,
		&task.CompletedAt,
		&task.TotalDuration,
	)

	if err != nil {
		return nil, err
	}
	return task, nil
}

//go:embed queries/status_task.sql
var statusTask string

//go:embed queries/start_task.sql
var startTask string

func StartTask(id, userID int) error {
	var status string

	err := client.Pool.QueryRow(client.Ctx, statusTask, id, userID).Scan(&status)
	if status == string(models.Started) {
		return fmt.Errorf("task already started, current status is '%s', err: %w", status, err)
	}

	count, err := client.Pool.Exec(client.Ctx, startTask, id, userID)
	if count.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return err
}

//go:embed queries/stop_task.sql
var stopTask string

func StopTask(id, userID int) (string, error) {
	var status string

	err := client.Pool.QueryRow(client.Ctx, statusTask, id, userID).Scan(&status)
	if status != string(models.Started) {
		return "", fmt.Errorf("task has not started, current status is '%s', err: %w", status, err)
	}

	var duration string

	err = client.Pool.QueryRow(client.Ctx, stopTask, id, userID).Scan(&duration)
	if err != nil {
		return "", err
	}

	return duration, nil
}

//go:embed queries/get_user_tasks.sql
var getUserTasks string

func GetUserTasks(userID int) ([]*models.Task, error) {
	var tasks []*models.Task

	rows, err := client.Pool.Query(client.Ctx, getUserTasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		err = rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.StartedAt,
			&task.CompletedAt,
			&task.TotalDuration,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
