package apimodels

type CreateCategoryRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
