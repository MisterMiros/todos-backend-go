package apimodels

import "go-todos/model"

type Todo struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
	CategoryId string `json:"category_id"`
}

func NewTodo(model *model.Todo) *Todo {
	return &Todo{
		Id: model.Id,
		Title: model.Title,
		Completed: model.Completed,
		CategoryId: model.CategoryId,
	}
}

func NewTodos(models []model.Todo) []Todo {
	result := []Todo{}
	for _, model := range models {
		result = append(result, *NewTodo(&model))
	}
	return result
}