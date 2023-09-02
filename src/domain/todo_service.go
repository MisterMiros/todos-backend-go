package domain

import (
	"fmt"
	"go-todos/model"
	"go-todos/storage"
)

type TodoService struct {
	categoryRepository *storage.CategoryRepository
	todoRepository     *storage.TodoRepository
}

func NewTodoService(todoRepository *storage.TodoRepository, categoryRepository *storage.CategoryRepository) *TodoService {
	return &TodoService{
		todoRepository:     todoRepository,
		categoryRepository: categoryRepository,
	}
}

func (service *TodoService) GetUserTodos(userEmail string) ([]model.Todo, *ServiceError) {
	todos, err := service.todoRepository.GetUserTodos(userEmail)
	if err != nil {
		return []model.Todo{}, NewInternalError(err)
	}
	return todos, nil
}

func (service *TodoService) GetTodo(userEmail, id string) (*model.Todo, *ServiceError) {
	todo, err := service.todoRepository.GetTodo(userEmail, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if todo == nil {
		return nil, NewNotFound()
	}
	return todo, nil
}

func (service *TodoService) CreateTodo(userEmail, title, categoryId string) (*model.Todo, *ServiceError) {
	categoryErr := service.doesCategoryExists(userEmail, categoryId)
	if (categoryErr != nil) {
		return nil, categoryErr
	}
	todo, err := service.todoRepository.CreateTodo(userEmail, title, categoryId)
	if err != nil {
		return nil, NewInternalError(err)
	}
	return todo, nil
}

func (service *TodoService) UpdateTodo(todo *model.Todo) *ServiceError {
	existingTodo, err := service.todoRepository.GetTodo(todo.UserEmail, todo.Id)
	if err != nil {
		return NewInternalError(err)
	}
	if existingTodo == nil {
		return NewNotFound()
	}

	categoryErr := service.doesCategoryExists(todo.UserEmail, todo.CategoryId)
	if (categoryErr != nil) {
		return categoryErr
	}

	err = service.todoRepository.UpdateTodo(todo)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func (service *TodoService) DeleteTodo(userEmail, id string) *ServiceError {
	err := service.todoRepository.DeleteTodo(userEmail, id)
	if err != nil {
		return NewInternalError(err)
	}
	return nil
}

func (service *TodoService) doesCategoryExists(userEmail, categoryId string) *ServiceError {
	if categoryId != "" {
		category, err := service.categoryRepository.GetCategory(userEmail, categoryId)
		if err != nil {
			return NewInternalError(err)
		}
		if category == nil {
			return NewBadRequest(fmt.Errorf("category '%v' not found", categoryId))
		}
	}
	return nil
}
