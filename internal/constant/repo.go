package constant

import (
	"errors"
)

// tables in DB
const (
	TasksTable    string = "tasks"
	StatusesTable string = "statuses"
)

// placeholder in sql query
const PlaceholderDollar string = "$"

// errors
var (
	ErrTaskIDNotExists = errors.New("no task with such an id")
)
