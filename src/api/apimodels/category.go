package apimodels

import "go-todos/model"

type Category struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func NewCategory(model *model.Category) *Category {
	return &Category{
		Id:    model.Id,
		Name:  model.Name,
		Color: model.Color,
	}
}

func NewCategories(models []model.Category) []Category {
	result := []Category{}
	for _, model := range models {
		result = append(result, *NewCategory(&model))
	}
	return result
}
