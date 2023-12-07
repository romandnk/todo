package service

type CreateStatusParams struct {
	Name string `json:"name" binding:"required"`
}

type CreateStatusResponse struct {
	ID int `json:"id"`
}
