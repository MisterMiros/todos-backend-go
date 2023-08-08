package main

import (
	"encoding/json"
	"go-todos/api/apimodels"
	"go-todos/api/utils"
	"go-todos/api/utils/responses"
	"go-todos/domain"
	"go-todos/model"
	"log"
	"strings"

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

	id, err := utils.GetId(event)
	if err != nil {
		log.Printf("Failed to get item 'id'. Error: %v", err)
		return responses.BadRequest("Failed to get item 'id'")
	}

	var body apimodels.UpdateTodoRequest
	err = json.Unmarshal([]byte(event.Body), &body)
	if err != nil {
		log.Printf("Failed to parse request body. Error: %v", err)
		return responses.InternalError("Failed to parse request body")
	}
	body.Title = strings.TrimSpace(body.Title)

	if len(body.Title) == 0 {
		return responses.BadRequest("Can't set empty title to an item")
	}

	todo := model.Todo{
		UserEmail: email,
		Id: id,
		Title: body.Title,
		Completed: body.Completed,
	}
	serviceErr := handler.todoService.UpdateTodo(&todo)
	if serviceErr != nil {
		log.Printf("Failed to create todo. Error: %v\n", serviceErr)
		return utils.ResponseFromServiceError(serviceErr)
	}
	return responses.NoContent()
}
