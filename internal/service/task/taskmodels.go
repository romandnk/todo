package taskservice

type CreateTaskParams struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	StatusName  string `json:"status_name" binding:"required"`
	Date        string `json:"date" binding:"required"`
}

type CreateTaskResponse struct {
	ID int `json:"id"`
}

type UpdateTaskByIDParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusName  string `json:"status_name"`
	Date        string `json:"date"`
}

type GetTaskWithStatusNameModel struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusName  string `json:"status_name"`
	Date        string `json:"date"`
	CreatedAt   string `json:"created_at"`
}

type GetAllTasksResponse struct {
	Total int                          `json:"total"`
	Tasks []GetTaskWithStatusNameModel `json:"tasks" json:"tasks"`
}
