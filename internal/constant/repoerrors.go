package constant

import "errors"

// task repo errors
var (
	ErrTaskIDNotExists = errors.New("no task with such an id")
)
