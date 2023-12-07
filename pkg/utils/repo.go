package utils

import (
	"database/sql"
	"fmt"
	"github.com/romandnk/todo/internal/entity"
	"strings"
)

// SetPlaceholders return string with defined quantity selected placeholders
func SetPlaceholders(placeholder string, quantity int) string {
	var result strings.Builder

	result.WriteString("(")

	for i := 1; i < quantity; i++ {
		result.WriteString(fmt.Sprintf("%s%d, ", placeholder, i))
	}

	result.WriteString(fmt.Sprintf("%s%d", placeholder, quantity))
	result.WriteString(")")

	return result.String()
}

type TaskToUpdate struct {
	Title       sql.NullString
	Description sql.NullString
	StatusID    sql.NullInt16
	Date        sql.NullTime
}

func CheckEmptyTaskFields(task entity.Task) TaskToUpdate {
	return TaskToUpdate{
		Title: sql.NullString{
			String: task.Title,
			Valid:  task.Title != "",
		},
		Description: sql.NullString{
			String: task.Description,
			Valid:  task.Description != "",
		},
		StatusID: sql.NullInt16{
			Int16: int16(task.StatusID),
			Valid: task.StatusID != 0,
		},
		Date: sql.NullTime{
			Time:  task.Date,
			Valid: !task.Date.IsZero(),
		},
	}
}
