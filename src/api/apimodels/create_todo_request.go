package apimodels

type CreateTodoRequest struct {
	Title      string `json:"title"`
	CategoryId string `json:"category_id"`
}