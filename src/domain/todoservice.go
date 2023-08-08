package domain

import (
	"go-todos/model"
	"go-todos/storage"
)

type TodoService struct {
	repository *storage.TodoRepository
}

func NewTodoService(repository *storage.TodoRepository) *TodoService {
	return &TodoService{
		repository: repository,
	}
}

func (service *TodoService) GetUserTodos(userEmail string) ([]model.Todo, *ServiceError) {
	todos, err := service.repository.GetUserTodos(userEmail)
	if err != nil {
		return []model.Todo{}, NewInternalError(err)
	}
	return todos, nil
}

func (service *TodoService) GetTodo(userEmail, id string) (*model.Todo, *ServiceError) {
	todo, err := service.repository.GetTodo(userEmail, id)
	if err != nil {
		return nil, NewInternalError(err)
	}
	if todo == nil {
		return nil, NewNotFound()
	}
	return todo, nil
}

func (service *TodoService) CreateTodo(userEmail, title string) (*model.Todo, *ServiceError) {
	todo, err := service.repository.CreateTodo(userEmail, title)
	if err != nil {
		return nil, NewInternalError(err)
	}
	return todo, nil
}

func (service *TodoService) UpdateTodo(todo *model.Todo) *ServiceError {
	existingTodo, err := service.repository.GetTodo(todo.UserEmail, todo.Id)
	if err != nil {
		return NewInternalError(err)
	}
	if existingTodo == nil {
		return NewNotFound()
	}

	err = service.repository.UpdateTodo(todo)
	if err != nil {
		return NewInternalError(err)
	}

	return nil
}

func (service *TodoService) DeleteTodo(userEmail, id string) *ServiceError {
	err := service.repository.DeleteTodo(userEmail, id)
	if err != nil {
		return NewInternalError(err)
	}
	return nil
}
