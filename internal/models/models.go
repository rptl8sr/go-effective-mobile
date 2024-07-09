package models

import "time"

type User struct {
	ID             int        `json:"id"`
	Surname        string     `json:"surname,omitempty"`
	Name           string     `json:"name,omitempty"`
	Patronymic     string     `json:"patronymic,omitempty"`
	Address        string     `json:"address,omitempty"`
	PassportNumber string     `json:"passport,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type Status string

var (
	Created = Status("created")
	Started = Status("started")
	Stopped = Status("stopped")
)

type Task struct {
	ID            int           `json:"id"`
	UserID        int           `json:"user_id"`
	Title         string        `json:"title"`
	Status        Status        `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     *time.Time    `json:"updated_at,omitempty"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
	TotalDuration time.Duration `json:"total_duration"`
}
