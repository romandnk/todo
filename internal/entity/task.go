package entity

import "time"

type Task struct {
	ID          int
	Title       string
	Description string
	StatusID    int
	Date        time.Time
	Deleted     bool
	CreatedAt   time.Time
	DeletedAt   time.Time
}
