package constant

import "errors"

// task repo errors
var (
	ErrTaskIDNotExists = errors.New("no task with id")
)

// utils repo errors
var (
	ErrNonPositiveQuantity = errors.New("quantity must be positive")
	ErrEmptyPlaceholder    = errors.New("placeholder cannot be empty")
)
