package apimodels

type UpdateTodoRequest struct {
	Title      string `json:"title"`
	CategoryId string `json:"category_id"`
	Completed  bool   `json:"completed"`
}
