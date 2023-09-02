package domain

import (
	"go-todos/model"
	"go-todos/storage"
)

type CategoryService struct {
	categoryRepository *storage.CategoryRepository
	todoRepository     *storage.TodoRepository
}

func NewCategoryService(categoryRepository *storage.CategoryRepository, todoRepository *storage.TodoRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
		todoRepository: todoRepository,
	}
}

func (service *CategoryService) GetUserCategories(userEmail string) ([]model.Category, *ServiceError) {
	categorys, err := service.categoryRepository.GetUserCategories(userEmail)
	if err != nil {
		return []model.Category{}, NewInternalError(err)
	}
	return categorys, nil
}

func (service *CategoryService) GetCategory(userEmail, id string) (*model.Category, *ServiceError) {
	category, err := service.categoryRepository.GetCategory(userEmail, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if category == nil {
		return nil, NewNotFound()
	}
	return category, nil
}

func (service *CategoryService) CreateCategory(userEmail, name, color string) (*model.Category, *ServiceError) {
	category, err := service.categoryRepository.CreateCategory(userEmail, name, color)
	if err != nil {
		return nil, NewInternalError(err)
	}
	return category, nil
}

func (service *CategoryService) UpdateCategory(category *model.Category) *ServiceError {
	existingCategory, err := service.categoryRepository.GetCategory(category.UserEmail, category.Id)
	if err != nil {
		return NewInternalError(err)
	}
	if existingCategory == nil {
		return NewNotFound()
	}

	err = service.categoryRepository.UpdateCategory(category)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func (service *CategoryService) DeleteCategory(userEmail, id string) *ServiceError {
	err := service.todoRepository.ClearCategory(userEmail, id)
	if err != nil {
		return NewInternalError(err)
	}
	err = service.categoryRepository.DeleteCategory(userEmail, id)
	if err != nil {
		return NewInternalError(err)
	}
	return nil
}
