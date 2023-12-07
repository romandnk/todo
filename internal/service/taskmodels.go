package service

type CreateTaskParams struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StatusName  string `json:"status_name" binding:"required"`
	Date        string `json:"date" binding:"required"`
}

type CreateTaskResponse struct {
	ID int `json:"id"`
}

type DeleteTaskParams struct {
	ID string `json:"id"`
}
