package constant

import "errors"

// mutual errors
var (
	ErrInternalError = errors.New("internal error")
)

// status service errors
var (
	ErrEmptyStatusName   = errors.New("status name cannot be empty")
	ErrTooLongStatusName = errors.New("max status name length is 36")
	ErrStatusNameExists  = errors.New("status name already exists")
)
