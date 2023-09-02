package domain

import (
	"go-todos/model"
	"go-todos/storage"
)

type CategoryService struct {
	repository *storage.CategoryRepository
}

func NewCategoryService(repository *storage.CategoryRepository) *CategoryService {
	return &CategoryService{
		repository: repository,
	}
}

func (service *CategoryService) GetUserCategories(userEmail string) ([]model.Category, *ServiceError) {
	categorys, err := service.repository.GetUserCategories(userEmail)
	if err != nil {
		return []model.Category{}, NewInternalError(err)
	}
	return categorys, nil
}

func (service *CategoryService) GetCategory(userEmail, id string) (*model.Category, *ServiceError) {
	category, err := service.repository.GetCategory(userEmail, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if category == nil {
		return nil, NewNotFound()
	}
	return category, nil
}

func (service *CategoryService) CreateCategory(userEmail, name, color string) (*model.Category, *ServiceError) {
	category, err := service.repository.CreateCategory(userEmail, name, color)
	if err != nil {
		return nil, NewInternalError(err)
	}
	return category, nil
}

func (service *CategoryService) UpdateCategory(category *model.Category) *ServiceError {
	existingCategory, err := service.repository.GetCategory(category.UserEmail, category.Id)
	if err != nil {
		return NewInternalError(err)
	}
	if existingCategory == nil {
		return NewNotFound()
	}

	err = service.repository.UpdateCategory(category)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func (service *CategoryService) DeleteCategory(userEmail, id string) *ServiceError {
	err := service.repository.DeleteCategory(userEmail, id)
	if err != nil {
		return NewInternalError(err)
	}
	return nil
}
