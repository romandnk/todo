package constant

import "errors"

// mutual errors
var (
	ErrInternalError = errors.New("internal error")
)

// status service errors
var (
	ErrEmptyStatusName   = errors.New("status name cannot be empty")
	ErrTooLongStatusName = errors.New("max status name length is 16")
	ErrStatusNameExists  = errors.New("status name already exists")
)

// task service errors
var (
	ErrEmptyTitle         = errors.New("title cannot be empty")
	ErrEmptyDescription   = errors.New("description cannot be empty")
	ErrTooLongTitle       = errors.New("max task title length is 64")
	ErrEmptyDate          = errors.New("date cannot be empty")
	ErrInvalidDateFormat  = errors.New("date must be in RFC3339 format")
	ErrOutdatedDate       = errors.New("you cannot set task date on the past")
	ErrEmptyTaskID        = errors.New("task id cannot be empty")
	ErrInvalidTaskID      = errors.New("task id must be int")
	ErrNonPositiveTaskID  = errors.New("task id must be positive")
	ErrInvalidLimit       = errors.New("limit must be int")
	ErrInvalidLastTaskID  = errors.New("last task id must be int")
	ErrNegativeLimit      = errors.New("limit cannot be negative")
	ErrNegativeLastTaskID = errors.New("last task id cannot be negative")
)
