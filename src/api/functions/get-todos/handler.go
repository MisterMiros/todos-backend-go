package main

import (
	"go-todos/api/utils"
	"go-todos/api/utils/responses"
	"go-todos/domain"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	todoService *domain.TodoService
}

func NewHandler(todoService *domain.TodoService) *Handler {
	return &Handler{
		todoService: todoService,
	}
}

func (handler *Handler) Handle(event events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	email, err := utils.GetEmail(event)
	if err != nil {
		log.Printf("Failed to parse authorizer context. Error: %v", err)
		return responses.BadRequest("Failed to parse authorizer context")
	}

	log.Printf("Getting list of todos for user '%v'\n", email)
	todos, serviceErr := handler.todoService.GetUserTodos(email)
	if serviceErr != nil {
		log.Printf("Failed to get user todos. Error: %v\n", serviceErr)
		return utils.ResponseFromServiceError(serviceErr)
	}
	return responses.Ok(todos)
}
